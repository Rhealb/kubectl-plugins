# kubectl-hostpathpv
　
## 说明

[English](README.md) | [中文](README-zh.md)

**Hostpath PV的绑定信息和使用情况多记录在对应PV的Annotations里面，在没有该工具之前想查看详细情况只能kubectl get pv $pvname -o json来查看，由于保存的是json数据查看起来不是那么方便．同时还有其他一些需要修改annotation以及数据迁移等高级操作，必须有一个类似kubectl的工具给用户使用所以才有了该工具的出现．目前该工具支持如下操作：**

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
 
+ **1) add命令：**

    为指定的PV在给定的Node上加一个使用记录，由于对于的hostpathpv插件不会根据该记录去创建目录所以不建议使用该命令．

+ **2) delete命令：**

    和add命令相反，为指定的PV删除Node上某一个使用记录(如果该目录没有Pod使用则该目录会在物理机上被删除数据丢失，所以只有在确认不需要该数据时才使用)．

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

+ **3) describe命令：**

    查看PV,Node,Pod的hostpathpv的详细使用情况．

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
     
+ **4) get命令：**

    查看使用hostpath pv的Pod,Node和PV.

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
    
+ **5) move命令：**

    将数据从一个node迁移到另一个node.
	
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

+ **6) scale命令：**

    对hostpathpv 进行扩容或者缩容．


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
     
+ **7) setdisable命令：**

    设置node的某个disk不可调度/可以调度，但是之前在上面创建的quota目录不受影响．

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
    
+ **8) upgrade命令：**

    刚开始实现hostpath的时候是接管的hostpath pv而且代码也是写在kubelet里面的，现在将其独立成CSI插件，该命令是将hostpath pv升级成CSI hostpathpv的．
	
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
    
## 部署
+ **1) 下载代码及编译：**

    	$ git clone ssh://git@gitlab.cloud.enndata.cn:10885/kubernetes/k8s-plugins.git
		$ cd k8s-plugins/kubectl-plugins/hostpathpv
		$ make install
    
    （make install 会将编译好的二进制文件kubectl-hostpathpv 移动到/user/bin目录下，从而可以通过使用kubectl hostpathpv 来间接使用kubectl-hostpathpv,也可以使用make build只生成kubectl-hostpathpv到当前目录，详情可看Makefile文件）

+ **2) 制作move所需image:**
	
    $ make moverelease REGISTRY=10.19.140.200:29006
    
    (该命令将会制作一个docker image 10.19.140.200:29006/library/hostpathscpmove:${TAG}并将其push到registry)
    
+ **3) 卸载：**

		$ make uninstall
