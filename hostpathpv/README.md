# kubectl-hostpathpv
　
## Explanation

[English](README.md) | [中文](README-zh.md)

**The binding information and usage of Hostpath PV are recorded in the Annotations of the corresponding PV，Using kubectl get pv $pvname -o yaml to view the details is not so intuitive．At the same time there are other advanced operations such as annotation modify and data migration  so we need a kubectl plugin to do these.**

	$ kubectl hostpathpv --help
      hostpathpv controls the Kubernetes cluster hostpathpv manager.

      Find more information at:
            https://gitlab.cloud.enndata.cn/kubernetes/k8s-plugins/kubectl-plugins/hostpathpv/

	Usage:
  		hostpathpv [flags]
  		hostpathpv [command]

	Available Commands:
  		add         Add pv quota path
  		delete      Delete pv quota path
  		describe    Describe node pod or pv 		quota information
  		get         Display node pod or pv quota information
  		help        Help about any command
  		move        Move quota path
  		scale       Scale up, down or to pv's capacity
  		setdisable  Disable or unDisable node quota disk
  		upgrade     Upgrade hostpath pv

	Flags:
  		-h, --help   help for hostpathpv

	Use "hostpathpv [command] --help" for more information about a command.
 
+ **1) add command：**

    Adding a usage record to a given Node for the specified PV, but it's not recommended because the hostpathpv plugin will not create a directory based on that record．

+ **2) delete command：**

    In contrast to the add command, delete a usage record on Node for the specified PV (if the directory is not used by Pod, the directory will be deleted on the physical machine and lost, so it will only be used if it is confirmed that the data is not needed)．

		$ kubectl hostpathpv delete --help
		Delete pv quota paths. 

		Valid resource types include: 

  			* pv

		Usage:
  			hostpathpv delete (TYPE/NAME ...) [flags]

		Examples:
	  	# Delete pv nodename quota path.
  			kubectl hostpathpv delete pv pvname --node=nodename --path=quotapath
  
  	  	# Delete pv all quota path.
  			kubectl hostpathpv delete pv pvname --all=true
  
  	  	# Delete pv nodename quota path without verify.
  			kubectl hostpathpv delete pv pvname --node=nodename --path=quotapath --force=true

		Flags:
      	--all           Delete all quota path
      	--force         Delete force with no verify
  	 	-h, --help          help for delete
      	--node string   Delete quota path of node
      	--path string   Delete path used with --node

+ **3) describe command：**

    View the details of hostpath PV, Node which supporting hostpath PV and Pod which used hostpath pv．

		$ kubectl hostpathpv describe --help
    	Show details of a specific resource diskquota info. 

		Valid resource types include: 
	
  			* node  
  			* pv  
  			* pod

		Usage:
  			hostpathpv describe (TYPE/NAME ...) [flags]

		Examples:
  			# Describe pods quota info.
  	 		kubectl hostpathpv describe pods
  
  			# Describe specified pod quota info.
 		 	kubectl hostpathpv describe pods podname
  
 			# Describe all node quota info.
  		 	kubectl hostpathpv describe nodes
  
  			# Describe a node quota info with specified NAME.
  	 		kubectl hostpathpv describe node nodename
  
  			# Describe all pvs quota inf.
   	 		kubectl hostpathpv describe pvs
  
  			# Describe a pv quota info with specified NAME.
  	 		kubectl hostpathpv describe pv pvname

		Flags:
			--all-namespaces     If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.
			--filter             Ignored no disk quota resource (default true)
 			 -h, --help               help for describe
 	 		-n, --namespace string   Namespace (default "default")
     
+ **4) get command：**

    View the Pod which used hostpathpv, Node which supporting hostpath PV and hostpath pv.

		$ kubectl hostpathpv get --help
     	Display one or many resources diskquota infomation. 

		Valid resource types include: 

 	 	* all  
 	 	* node  
 	 	* pv  
 	 	* pod

		Usage:
  			hostpathpv get (TYPE/NAME ...) [flags]

		Examples:
  			# List all pods quota info in ps output format.
  			kubectl hostpathpv get pods
  
  			# List specified pod quota info in ps output format.
  			kubectl hostpathpv get pods podname
  
  			# List all node quota info in ps output format.
  			kubectl hostpathpv get nodes
  
  			# List a node quota info with specified NAME in ps output format.
  			kubectl hostpathpv get node nodename
  
  			# List all pvs quota info in ps output format.
 		 	kubectl hostpathpv get pvs
  
  			# List a pv quota info with specified NAME in ps output format.
  			kubectl hostpathpv get pv pvname
  
  			# List all resources quota info.
 		 	kubectl hostpathpv get all

		Flags:
      		--all-namespaces     If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.
      		--filter             Ignored no disk quota resource (default true)
  			-h, --help               help for get
  			-n, --namespace string   namespace (default "default")
    
+ **5) move command：**

    Migrating data from one node to another.
	
    	$ kubectl hostpathpv move --help
    	Move node quota path from one disk to other disk. 

		Valid resource types include: 

 	 	* node

		Usage:
  			hostpathpv move (TYPE/NAME ...) [flags]

		Examples:
 			# Move quota path.
 			kubectl hostpathpv move node nodename --from=/xfs/disk1/dir --to=/xfs/disk2/
  
 			kubectl hostpathpv move node --from=node1:/xfs/disk1/dir1,/xfs/disk2/dir2 --to=node2:/xfs/disk2,/xfs/disk3

		Flags:
      		--alwayspullmoveimage   Move pod pull images alwasy
      		--force                 Move force with no confirmation
     		 --from string           Move quota path from
 	 		-h, --help                  help for move
     		 --moveimage string      Image create to move dir (default "127.0.0.1:29006/library/hostpathscpmove:v3.1")
     		 --movepodmemlimit int   Move pod memory MB (default 1024)
     		 --movetimeout int       Move pod scp timeout (default 100000)
     		 --tmppvkeepwait int     Wait time at create tmp pv (default 60)
     		 --to string             Move quota to

+ **6) scale command：**

    Expansion or contraction of hostpath PV．


		$ kubectl hostpathpv scale --help
   		 Scale capacity. 

		 Valid resource types include: 

  		* node

		Usage:
 			hostpathpv scale (TYPE/NAME ...) [flags]

		Examples:
  			# Scale up pv's capacity example.
  			kubectl hostpathpv scale pv pvname up 100
  
  			# Scale down pv's capacity example.
  			kubectl hostpathpv scale pv pvname down 100
  
  			# Scale pv's capacity to example.
  			kubectl hostpathpv scale pv pvname to 1100

		Flags:
 	 		-h, --help   help for scale
     
+ **7) setdisable command：**

    Setting a disk of node is not schedulable / schedulable, but the quota directory previously created above is unaffected．

		$ kubectl hostpathpv setdisable --help
		Disable or unDisable node quota disk. 

		Valid resource types include: 

		* node

		Usage:
			hostpathpv setdisable (TYPE/NAME ...) [flags]
        
		Examples:
			# Disable node nodename quota disk.
  			kubectl hostpathpv setdisable node nodename --disk diskpath --disable=true
  
 			# Undisable node nodename quota disk.
  			kubectl hostpathpv setdisable node nodename --disk diskpath --disable=false

		Flags:
    		  --disable       Set disk diskable or undiskable (default true)
    	  --disk string   Disable or unDisable quota disk of node
  			-h, --help          help for setdisable
    
+ **8) upgrade command：**

    Upgrades hostpath PV to CSI hostpath pv．
	
    	$ kubectl hostpathpv upgrade --help
    	Upgrade hostpath pv to CSI hostpathpv. 

		Valid resource types include: 

  			* pv

		Usage:
 	 	 hostpathpv upgrade (TYPE/NAME ...) [flags]

		Examples:
  			# Upgrade quota path.
  			kubectl hostpathpv upgrade pv pvname

		Flags:
     	 --deleteinterval duration   Use quota pod deleting interval (default 10s)
     	 --force                     Upgrade force
  		-h, --help                      help for upgrade
     	 --upgradeimage string       Image create to change quota dir type (default "127.0.0.1:29006/library/busybox:1.25")
    
## Deploy
+ **1) Code download  and compile：**

    	$ git clone ssh://git@gitlab.cloud.enndata.cn:10885/kubernetes/k8s-plugins.git
		$ cd k8s-plugins/kubectl-plugins/hostpathpv
		$ make install
    
    （Make install will move the compiled binary file kubectl-hostpathpv to the / user/bin directory, so you can use 'kubectl hostpathpv' indirectly by using kubectl-hostpathpv, or you can use make build to generate only kubectl-hostpathpv to the current directory）

+ **2) build move command needed image:**
	
    $ make moverelease REGISTRY=10.19.140.200:29006
    
    (This command will make a docker image 10.19.140.200:29006/library/hostpath scpmove:${TAG} and push it to registry)
    
+ **3) Uninstall：**

		$ make uninstall
