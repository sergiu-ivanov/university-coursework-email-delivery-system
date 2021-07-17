# Email Delivery System

An email delivery system that implements micro service architecture in Go.

<img width="1210" alt="Screenshot 2021-07-17 at 19 47 15" src="https://user-images.githubusercontent.com/43847681/126044173-4d4f113a-8832-4649-b172-7e1feeae79f3.png">


# INSTRUCTIONS:

1. BUILD an image with the following command:				
								
docker build --tag [IMAGE NAME] .					
							

2. CREATE a virtual network with the following command:		
								
docker network create --subnet [IP ADDRESS] [NETWORK NAME]		


3. ADD a container container to the network with the following command:
								
docker run --name [CONTAINER NAME] --net [NETWORK NAME] \		
              --ip [MSA's IP] --detach \				
              --publish [UNALLOCATED PORT]:8888 \			
              --security-opt apparmor=unconfined [IMAGE NAME]		

