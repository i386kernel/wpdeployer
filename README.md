# Wordpress Auto Deploy
## This repo contains utilities to auto-deploy wordpress and enable resiliency

[Medium Link](https://i386kernel.medium.com/how-to-make-cloud-native-work-loads-truly-resilient-to-disasters-ba6f0fe15aa7)

### Running wphelper service
```sh
[root@NFS]# systemctl status wphelper.service
● wphelper.service - Wordpress YAML Generator
   Loaded: loaded (/etc/systemd/system/wphelper.service; enabled; vendor preset: disabled)
   Active: active (running) since Fri 2020-11-27 11:58:53 IST; 1 weeks 3 days ago
 Main PID: 19415 (wphelper)
    Tasks: 6
   CGroup: /system.slice/wphelper.service
           └─19415 /bin/wphelper`
```

### How to use wpdeployer
```sh
[root@rhel wpdeployer]# ./wpdeployer
------------Performing NFS Operations----------
Checking if the project exists
Proceeding with a new ---wp-auto-31763--- project
Performing NFS Directory Operations
Performing NFS Service Operations
Performing Wordpress helper service operations
---------Performing OCP operations----------
Creating Project in <Openshift URL>
---Project Creation Status: 201---
MySql Operations
Creating mysql Persistent Volume
Creating PV mysql-wp-auto-31763
---PV Status: 201---
Creating Mysql Persistent volume claim
Creating PVC mysql-wp-auto-31763
---PVC Status: 201---
Creating Mysql Deployment
Deploying Mysql mysql-wp-auto-31763
---MySql Deployment status: 201---
Creating Mysql Service
Creating MySql Service mysql-wp-auto-31763
---MySQL Service Status: 201---
Creating wordpress Persistent Volume
Creating PV wordpress-wp-auto-31763
---PV Status: 201---
Creating wordpress Persistent volume Volume
Creating PVC wordpress-wp-auto-31763
---PVC Status: 201---
Creating wordpress deployment
Deploying Wordpress wordpress-wp-auto-31763
---Wordpress Deployment status: 201---
Creating wordpress service
Creating Wordpress Service wordpress-wp-auto-31763
---Wordpress Service Status: 201---
Performing DR Volume Operations
Creating MySql Persistent Volume
Creating PV mysql-wp-auto-31763
---PV Status: 201---
Creating Wordpress Persistent Volume
Creating PV wordpress-wp-auto-31763
---PV Status: 201---
Time Elapsed: 4.083705153s
Wordpress PR URL: <PR SVC>
Wordpress DR URL: <DR SVC>
```
