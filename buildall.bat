@echo on
cd admin
call pnpm i
call pnpm run build
cd ..

cd cli
call pnpm i
call pnpm run build
cd ..

cd ext

cd control
call pnpm i 
call pnpm run build
cd ..

cd ai
call pnpm i 
call pnpm run build
cd ..

cd pgsql
call pnpm i 
call pnpm run build
cd ..

cd mysql
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

cd sdk/js
call pnpm i
call pnpm update
call pnpm run build
cd ../..

go build