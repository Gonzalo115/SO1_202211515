while true;
do
    docker ps -a | awk 'NR>1{print $1}' | xargs sudo docker rm -f
    sleep 10
done