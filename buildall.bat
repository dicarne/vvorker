@echo on
cd www
call pnpm i
call pnpm run build
call pnpm run export
cd ..

cd ext/ai
call pnpm i 
call pnpm run build
cd ..

cd ..

go build