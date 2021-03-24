
curl_get() {
  local times=$1
  local id=$2

  for ((i=1;i<=$times;i++))
  do
    curl -X GET http://localhost:3000/user?id=$id
    echo "------- request GET / ${i} times --------- "
  done
}

curl_post() {
  local times=4
  response=$(curl -X POST -H "Accept: application/json" -d '{"name" : "Zoe"}' http://localhost:3000/user | jq -r '.userID') 
  echo "User id is $response"
  echo "get user name $times"
  
  curl_get $times $response
}

curl_post 
