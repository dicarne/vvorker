@echo on
cd www
call pnpm i
call pnpm run build
call pnpm run export
cd ..

cd ext

cd ai
call pnpm i 
call pnpm run build
cd ..

cd pgsql
call pnpm i 
call pnpm run build
cd ..

cd oss
call pnpm i 
call pnpm run build
cd ..

cd kv
call pnpm i 
call pnpm run build
cd ..

cd assets
call pnpm i 
call pnpm run build
cd ..

cd task
call pnpm i 
call pnpm run build
cd ..

cd ..
go build