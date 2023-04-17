clear
docker image build -t forum . 
docker run -it -p 4000:4000 forum