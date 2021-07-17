# Email Delivery System

An email delivery system that implements micro service architecture in Go.



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

