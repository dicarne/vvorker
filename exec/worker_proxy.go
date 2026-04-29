package exec

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
	"vvorker/defs"
	workercopy "vvorker/models/worker_copy"

	"github.com/sirupsen/logrus"
)

const workerDrainTimeout = 60 * time.Second

type workerProcessInfo struct {
	WorkerUID string
	CopyKey   string
}

type proxyTarget struct {
	address           string
	activeConnections atomic.Int64
}

func (t *proxyTarget) isIdle() bool {
	if t == nil {
		return true
	}
	return t.activeConnections.Load() == 0
}

type switchableTCPProxy struct {
	listenPort uint
	listener   net.Listener
	mu         sync.RWMutex
	current    *proxyTarget
}

func newSwitchableTCPProxy(listenPort uint) (*switchableTCPProxy, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", defs.DefaultHostName, listenPort))
	if err != nil {
		return nil, err
	}

	proxy := &switchableTCPProxy{
		listenPort: listenPort,
		listener:   listener,
	}

	go proxy.serve()
	return proxy, nil
}

func (p *switchableTCPProxy) serve() {
	for {
		clientConn, err := p.listener.Accept()
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
				continue
			}
			logrus.WithError(err).Warnf("worker proxy listener on %d stopped", p.listenPort)
			return
		}

		go p.handle(clientConn)
	}
}

func (p *switchableTCPProxy) handle(clientConn net.Conn) {
	target := p.currentTarget()
	if target == nil {
		_ = clientConn.Close()
		return
	}

	upstreamConn, err := net.Dial("tcp", target.address)
	if err != nil {
		logrus.WithError(err).Warnf("worker proxy dial failed, backend=%s", target.address)
		_ = clientConn.Close()
		return
	}

	target.activeConnections.Add(1)
	defer target.activeConnections.Add(-1)

	proxyTCP(clientConn, upstreamConn)
}

func (p *switchableTCPProxy) currentTarget() *proxyTarget {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.current
}

func (p *switchableTCPProxy) SwitchTo(address string) *proxyTarget {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.current != nil && p.current.address == address {
		return nil
	}

	oldTarget := p.current
	p.current = &proxyTarget{address: address}
	return oldTarget
}

func proxyTCP(clientConn net.Conn, upstreamConn net.Conn) {
	errChan := make(chan error, 1)
	onClose := func(_ error) {
		_ = upstreamConn.Close()
		_ = clientConn.Close()
	}

	go func() {
		_, err := io.Copy(clientConn, upstreamConn)
		errChan <- err
		onClose(err)
	}()

	go func() {
		_, err := io.Copy(upstreamConn, clientConn)
		errChan <- err
		onClose(err)
	}()

	<-errChan
}

func buildCopyKey(workerUID string, localID uint) string {
	return fmt.Sprintf("%s-%d", workerUID, localID)
}

func buildProcessKey(copy *workercopy.WorkerCopy, version string) string {
	backendPort := copy.BackendPort
	if backendPort == 0 {
		backendPort = copy.Port
	}
	return fmt.Sprintf("%s@%s@%d", buildCopyKey(copy.WorkerUID, copy.LocalID), version, backendPort)
}

func waitForTCPPort(port uint, timeout time.Duration) bool {
	address := fmt.Sprintf("%s:%d", defs.DefaultHostName, port)
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", address, 500*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return true
		}
		time.Sleep(200 * time.Millisecond)
	}
	return false
}

func (m *execManager) ensureProxy(key string, port uint) (*switchableTCPProxy, error) {
	if proxy, ok := m.proxyMap.Get(key); ok {
		return proxy, nil
	}

	proxy, err := newSwitchableTCPProxy(port)
	if err != nil {
		return nil, err
	}

	if existingProxy, ok := m.proxyMap.Get(key); ok {
		_ = proxy.listener.Close()
		return existingProxy, nil
	}

	m.proxyMap.Set(key, proxy)
	return proxy, nil
}

func (m *execManager) requestProcessStop(processKey string) {
	m.signMap.Set(processKey, true)
	if channel, ok := m.chanMap.Get(processKey); ok {
		select {
		case channel <- struct{}{}:
		default:
		}
	}
}

func (m *execManager) forceStopProcess(processKey string) {
	m.requestProcessStop(processKey)

	pid, ok := m.pidMap.Get(processKey)
	if !ok {
		return
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return
	}

	_ = process.Kill()
}

func (m *execManager) retireProcessAfterDrain(processKey string, mainTarget *proxyTarget, controlTarget *proxyTarget) {
	timeout := time.NewTimer(workerDrainTimeout)
	ticker := time.NewTicker(200 * time.Millisecond)
	defer timeout.Stop()
	defer ticker.Stop()

	for {
		if mainTarget.isIdle() && controlTarget.isIdle() {
			m.requestProcessStop(processKey)
			return
		}

		select {
		case <-timeout.C:
			logrus.Warnf("worker process %s drain timeout reached, forcing stop", processKey)
			m.forceStopProcess(processKey)
			return
		case <-ticker.C:
		}
	}
}

func (m *execManager) refreshWorkerRunning(workerUID string) {
	isRunning := false
	m.currentProcessMap.Range(func(_ string, processKey string) bool {
		info, ok := m.processInfoMap.Get(processKey)
		if !ok || info.WorkerUID != workerUID {
			return true
		}

		if _, ok := m.chanMap.Get(processKey); ok {
			isRunning = true
			return false
		}

		return true
	})
	m.runningMap.Set(workerUID, isRunning)
}

func (m *execManager) DrainAndExitCopy(workerUID string, localID uint) {
	copyKey := buildCopyKey(workerUID, localID)
	processKey, ok := m.currentProcessMap.Get(copyKey)
	if !ok {
		return
	}

	m.currentProcessMap.Delete(copyKey)
	m.refreshWorkerRunning(workerUID)
	m.signMap.Set(processKey, true)

	var mainTarget *proxyTarget
	if mainProxy, ok := m.proxyMap.Get(copyKey); ok {
		mainTarget = mainProxy.currentTarget()
	}

	var controlTarget *proxyTarget
	if controlProxy, ok := m.proxyMap.Get(copyKey + "-control"); ok {
		controlTarget = controlProxy.currentTarget()
	}

	go m.retireProcessAfterDrain(processKey, mainTarget, controlTarget)
}
