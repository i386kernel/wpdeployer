package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func ocpinteract(baseurl, endpoint, token, jsonbody string) (int, []byte){

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := http.Client{Transport: tr}
	request, err := http.NewRequest("POST", baseurl+endpoint, bytes.NewBuffer([]byte(jsonbody)))
	if err != nil {
		fmt.Println(err)
	}
	request.Header = map[string][]string{"Content-type": {"application/json"}, "Authorization": {"Bearer " + token}}

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	if resp.StatusCode > 210 {
		fmt.Println(string(body))
	}
	return resp.StatusCode, body
}

func ocppatch(baseurl, endpoint, token, jsonbody string) (int, []byte){

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := http.Client{Transport: tr}
	request, err := http.NewRequest("PATCH", baseurl+endpoint, bytes.NewBuffer([]byte(jsonbody)))
	if err != nil {
		fmt.Println(err)
	}
	request.Header = map[string][]string{"Content-type": {"application/strategic-merge-patch+json"}, "Authorization": {"Bearer " + token}}

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
	return resp.StatusCode, body
}

func createProject(baseurl, token string) string{
	fmt.Printf("Creating Project in '%s' \n", baseurl)
	var project string = "/apis/project.openshift.io/v1/projects/"
	jsonstr := fmt.Sprintf(`{
  "apiVersion": "project.openshift.io/v1",
  "kind": "Project",
  "metadata": {
    "annotations": {
      "openshift.io/requester": "%s"
    },
    "name": "%s"
     },
  "spec": {
    "finalizers": [
      "kubernetes"
    ]
  }
}`, "rouser", projectname)
	status, _ := ocpinteract(baseurl, project, token, jsonstr)
	fmt.Printf("---Project Creation Status: %d---\n", status)
	if status != 201{
		fmt.Println("Un Authorized, check your URL and TOKENS...Exiting..")
		os.Exit(1)
	}
	return projectname
}

func patchdeployment(url, token string, replicas int8) {
	deployment := fmt.Sprintf("/apis/apps/v1/namespaces/%s/deployments/wordpress-deployment", projectname)
	jsonstr := fmt.Sprintf(`{
	"spec": {
		"replicas": %d
	}
}`, replicas)
	status, _ := ocppatch(url, deployment, token, jsonstr)
	fmt.Println(status)
}

func createDeployment(url, token string) {
	deployment := fmt.Sprintf("/apis/apps/v1/namespaces/%s/deployments/", projectname)
	jsonstr := fmt.Sprintf(`{
  "apiVersion": "apps/v1",  
  "kind": "Deployment",
  "metadata": {
    "name": "wordpress-deployment"
  },
  "spec": {
	"replicas": 1,
	"serviceAccountName":"runasanyuid",
    "selector": {
      "matchLabels": {
        "app": "%s",
        "tier": "frontend"
      }
    },
    "strategy": {
      "type": "Recreate"
    },
    "template": {
      "metadata": {
        "labels": {
          "app": "%s",
          "tier": "frontend"
        }
      },
      "spec": {
        "containers": [
          {
            "name": "wordpress-mysql",
            "image": "mysql:5.6",
            "env": [
              {
                "name": "MYSQL_ROOT_PASSWORD",
                "value": "password"
              },
              {
                "name": "WORDPRESS_DB_NAME",
                "value": "wordpress"
              }
            ],
            "ports": [
              {
                "name": "wordpress-mysql",
                "containerPort": 3306
              }
            ],
            "volumeMounts": [
              {
                "name": "wordpress-persistent-storage",
                "mountPath": "/var/lib/mysql"
              }
            ]
          },
          {
            "name": "wordpress",
            "image": "192.168.8.201/resiliencyorchestration/wordpress:4.8-apache",
            "env": [
              {
                "name": "WORDPRESS_DB_HOST",
                "value": "127.0.0.1:3306"
              },
              {
                "name": "WORDPRESS_DB_USER",
                "value": "root"
              },
              {
                "name": "WORDPRESS_DB_PASSWORD",
                "value": "password"
              },
              {
                "name": "WORDPRESS_DB_NAME",
                "value": "wordpress"
              }
            ],
            "volumeMounts": [
              {
                "name": "wordpress-persistent-storage",
                "mountPath": "/var/www/html"
              }
            ]
          },
          {
            "name": "wordpress-db-updater",
            "image": "192.168.8.201/resiliencyorchestration/wordpress-db-updater:2.0",
            "volumeMounts": [
              {
                "name": "wordpress-persistent-storage",
                "mountPath": "/dbfile/"
              }
            ],
            "ports": [
              {
                "name": "wordpress",
                "containerPort": 80
              }
            ]
          }
        ],
        "volumes": [
          {
            "name": "wordpress-persistent-storage",
            "persistentVolumeClaim": {
              "claimName": "%s"
            }
          }
        ]
      }
    }
  }
}`, projectname, projectname, projectname+"-pvc")
	status, _ := ocpinteract(url, deployment, token, jsonstr)
	fmt.Println(status)
}

func (o *Openshift) createMySqlDep(){
	fmt.Printf("Deploying Mysql %s\n", o.Workloadname)
	deployment := fmt.Sprintf("/apis/apps/v1/namespaces/%s/deployments/", projectname)
	jsonstr := fmt.Sprintf(`{
  "apiVersion": "apps/v1",
  "kind": "Deployment",
  "metadata": {
    "name": "%s-mysql-wordpressmysql-deployment"
  },
  "spec": {
	"securityContext": {
				"fsGroup": 0
					},
    "selector": {
      "matchLabels": {
        "app": "wordpress",
        "tier": "backend"
      }
    },
    "strategy": {
      "type": "Recreate"
    },
    "template": {
      "metadata": {
        "labels": {
          "app": "wordpress",
          "tier": "backend"
        }
      },
      "spec": {
        "containers": [
          {
            "name": "wordpress-mysql",
            "image": "registry.access.redhat.com/rhscl/mysql-57-rhel7",
            "env": [
              {
                "name": "MYSQL_ROOT_PASSWORD",
                "value": "password"
              },
              {
                "name": "MYSQL_USER",
                "value": "wpuser"
              },
              {
                "name": "MYSQL_PASSWORD",
                "value": "password"
              },
              {
                "name": "WORDPRESS_DB_NAME",
                "value": "wordpress"
              }
            ],
			"securityContext": {
				"fsGroup": 0 
					},
            "ports": [
              {
                "name": "wordpress-mysql",
                "containerPort": 3306
              }
            ],
            "volumeMounts": [
              {
                "name": "wordpress-persistent-storage",
                "mountPath": "/var/lib/mysql"
              }
            ]
          }
        ],
        "volumes": [
          {
            "name": "wordpress-persistent-storage",
            "persistentVolumeClaim": {
              "claimName": "%s-pvc"
            }
          }
        ]
      }
    }
  }
}`, o.Workloadname, o.Workloadname)
	status, _ := ocpinteract(o.URL, deployment, o.Token, jsonstr)
	fmt.Printf("---MySql Deployment status: %d---\n", status)
}

func (o *Openshift) createWPDep(){
	fmt.Printf("Deploying Wordpress %s\n", o.Workloadname)
	deployment := fmt.Sprintf("/apis/apps/v1/namespaces/%s/deployments/", projectname)
	jsonstr := fmt.Sprintf(`{
   "apiVersion": "apps/v1",
   "kind": "Deployment",
   "metadata": {
      "name": "%s-wordpress-deployment"
   },
   "spec": {
      "selector": {
         "matchLabels": {
            "app": "wordpress",
            "tier": "frontend"
         }
      },
      "strategy": {
         "type": "Recreate"
      },
      "template": {
         "metadata": {
            "labels": {
               "app": "wordpress",
               "tier": "frontend"
            }
         },
         "spec": {
            "containers": [
               {
                  "name": "wordpress",
                  "image": "wordpress:php7.4-apache",
                  "env": [
                     {
                        "name": "WORDPRESS_DB_HOST",
                        "value": "mysqlsvc"
                     },
                     {
                        "name": "WORDPRESS_DB_USER",
                        "value": "root"
                     },
                     {
                        "name": "WORDPRESS_DB_PASSWORD",
                        "value": "password"
                     },
                     {
                        "name": "WORDPRESS_DB_NAME",
                        "value": "wordpress"
                     }
                  ],
                  "volumeMounts": [
                     {
                        "name": "wordpress-persistent-storage",
                        "mountPath": "/var/www/html"
                     }
                  ],
                  "ports": [
                     {
                        "name": "wordpress",
                        "containerPort": 80
                     }
                  ]
               },
				{
                  "name": "wordpress-db-updater",
                  "image": "docker.io/i386kernel/wpmanager:8.0",
                  "volumeMounts": [
                     {
                        "name": "wordpress-persistent-storage",
                        "mountPath": "/dbfile/"
                     }
                  ]
               }
            ],
            "volumes": [
               {
                  "name": "wordpress-persistent-storage",
                  "persistentVolumeClaim": {
                     "claimName": "%s-pvc"
                  }
               }
            ]
         }
      }
   }
}`, o.Workloadname, o.Workloadname)
	status, _ := ocpinteract(o.URL, deployment, o.Token, jsonstr)
	fmt.Printf("---Wordpress Deployment status: %d---\n", status)
}

func (o *Openshift)createPersistantVolume() {
	fmt.Printf("Creating PV %s\n", o.Workloadname)
	persistantVolume := "/api/v1/persistentvolumes/"
	jsonstr := fmt.Sprintf(`{
  "apiVersion": "v1",
  "kind": "PersistentVolume",
  "metadata": {
    "name": "%s-pv"
  },
  "spec": {
    "capacity": {
      "storage": "5Gi"
    },
    "volumeMode": "Filesystem",
    "accessModes": [
      "ReadWriteMany"
    ],
    "persistentVolumeReclaimPolicy": "Recycle",
    "storageClassName": "%s-storageclass",
    "mountOptions": [
      "hard",
      "nfsvers=4.2"
    ],
    "nfs": {
      "path": "/mnt/%s",
      "server": "%s"
    }
  }
}`, o.Workloadname, o.Workloadname, o.DIRName, NFS_IP)
	status, _ := ocpinteract(o.URL, persistantVolume, o.Token, jsonstr)
	fmt.Printf("---PV Status: %d---\n", status)
}


func (o *Openshift)createPersistantVolumeClaim() {
	fmt.Printf("Creating PVC %s\n", o.Workloadname)
	persistantVolumeClaim := fmt.Sprintf("/api/v1/namespaces/%s/persistentvolumeclaims/", projectname)
	jsonstr := fmt.Sprintf(`{
  "apiVersion": "v1",
  "kind": "PersistentVolumeClaim",
  "metadata": {
    "name": "%s-pvc",
    "labels": {
      "name": "wordpress-app-pv"
    }
  },
  "spec": {
    "accessModes": [
      "ReadWriteMany"
    ],
    "storageClassName": "%s-storageclass",
    "resources": {
      "requests": {
        "storage": "5Gi"
      }
    }
  }
}`, o.Workloadname, o.Workloadname)
	status, _ := ocpinteract(o.URL, persistantVolumeClaim, o.Token, jsonstr)
	fmt.Printf("---PVC Status: %d---\n", status)
}

func (o *Openshift)createWPService() {
	fmt.Printf("Creating Wordpress Service %s\n", o.Workloadname)
	service := fmt.Sprintf("/api/v1/namespaces/%s/services/", projectname)
	jsonstr := fmt.Sprintf(`{
  "apiVersion": "v1",
  "kind": "Service",
  "metadata": {
    "name": "wordpress-svc",
    "labels": {
      "app": "wordpress",
      "tier": "frontend"
    }
  },
  "spec": {
    "type": "NodePort",
    "ports": [
      {
        "name": "wordpres-pod-port",
        "port": 80,
        "targerPort": 80,
        "nodePort": %d
      }
    ],
    "selector": {
      "app": "wordpress",
      "tier": "frontend"
    }
  }
}`, randint)
	// nodeport range 30000-32767
	status, _ := ocpinteract(o.URL, service, o.Token, jsonstr)
	fmt.Printf("---Wordpress Service Status: %d---\n", status)
}

func (o *Openshift)createMySqlService() {
	fmt.Printf("Creating MySql Service %s\n", o.Workloadname)
	service := fmt.Sprintf("/api/v1/namespaces/%s/services/", projectname)
	jsonstr := fmt.Sprintf(`{
  "apiVersion": "v1",
  "kind": "Service",
  "metadata": {
    "name": "mysqlsvc",
    "labels": {
      "app": "wordpress",
      "tier": "backend"
    }
  },
  "spec": {
    "ports": [
      {
        "name": "mysql-port",
        "port": 3306
      }
    ],
    "selector": {
      "app": "wordpress",
      "tier": "backend"
    },
	"clusterIP": "None"
  }
}`)
	status, _ := ocpinteract(o.URL, service, o.Token, jsonstr)
	fmt.Printf("---MySQL Service Status: %d---\n", status)
}