# 解决方案
1. 创建两个镜像:1. 通过指定pvcname 备份leveldb 数据到s3； 2. 恢复s3 备份数据到指定pvc (针对levelDB 数据，通过创建snapshot 来读取并复制数据)
2. 创建job manager，通过传入参数调用创建指定镜像job 完成备份数据和恢复数据

# hints:
同个pvc 可以挂多个pod 只要pod 在同个node; mount S3 to kubernetes pod （https://dev.to/otomato_io/mount-s3-objects-to-kubernetes-pods-12f5）