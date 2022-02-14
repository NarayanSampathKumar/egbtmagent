# egbtmagent
<h1>Steps to run the operator:</h1>
<h3>Step - 1: Clone the code from github</h3>
git clone https://github.com/NarayanSampathKumar/egbtmagent.git  <br>
cd egbtmagent 
<h3>Step - 2: Make the files executable:</h3>
   chmod 777 redeploy_operator.sh bin/* <br>
   
<h3>Step - 3. change the repository from 172.16.8.78:5000 to egapm in the</h3>
   1.	Makefile<br>
      “IMAGE_TAG_BASE ?= 172.16.8.78:5000/egbtmagent” to “IMAGE_TAG_BASE ?= egapm/egbtmagent”<br>
   2.	And redeploy_operator.sh <br>
      “operator-sdk run bundle 172.16.8.78:5000/egbtmagent-bundle:v0.0.1 --skip-tls --timeout 15m0s” to “operator-sdk run bundle egapm/egbtmagent-bundle:v0.0.1 --timeout 15m0s”<br>
<h3>Step - 4. Execute the sh file with ./redeploy_operator.sh command</h3>
<h3>Step - 5. To cleanup the operator Run:</h3>
      operator-sdk cleanup egbtmagent --cleanup-all=true<br>
<h3>Step - 6. If any error during "operator-sdk run bundle" command run</h3>
     "operator-sdk cleanup egbtmagent --cleanup-all=true"<br>
   and run command<br>
     "operator-sdk run bundle egapm/egbtmagent-bundle:v0.0.1 --timeout 15m0s" <br>
