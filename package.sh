# compile for version
make
if [ $? -ne 0 ]; then
    echo "make error"
    exit 1
fi 

# cross_compiles
make -f ./Makefile.cross-compiles

rm -rf ./release/packages
mkdir -p ./release/packages

os_all='linux windows darwin freebsd'
arch_all='386 amd64 arm arm64 mips64 mips64le mips mipsle riscv64'

cd ./release

for os in $os_all; do
    for arch in $arch_all; do
        nbtverify_dir_name="nbtverify_${nbtverify_version}_${os}_${arch}"
        nbtverify_path="./packages/nbtverify_${nbtverify_version}_${os}_${arch}"

        if [ "x${os}" = x"windows" ]; then
            if [ ! -f "./nbtverify_${os}_${arch}.exe" ]; then
                continue
            fi 
            mkdir ${nbtverify_path}
            mv ./nbtverify_${os}_${arch}.exe ${nbtverify_path}/nbtverify.exe
        else
            if [ ! -f "./nbtverify_${os}_${arch}" ]; then
                continue
            fi 
            mkdir ${nbtverify_path}
            mv ./nbtverify_${os}_${arch} ${nbtverify_path}/nbtverify
        fi  
        cp ../LICENSE ${nbtverify_path}
        cp -rf ../conf/* ${nbtverify_path}

        # packages
        cd ./packages
        if [ "x${os}" = x"windows" ]; then
            zip -rq ${nbtverify_dir_name}.zip ${nbtverify_dir_name}
        else
            tar -zcf ${nbtverify_dir_name}.tar.gz ${nbtverify_dir_name}
        fi  
        cd ..
        rm -rf ${nbtverify_path}
    done
done

cd -
