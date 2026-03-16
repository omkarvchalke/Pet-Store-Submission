#!/bin/bash

URL="http://localhost:8080/query"
AUTH="merchant_happy:merchant123"

create_pet() {
  local name="$1"
  local species="$2"
  local age="$3"
  local description="$4"

  curl -s -u "$AUTH" \
    -X POST "$URL" \
    -H "Content-Type: application/json" \
    -d "{
      \"query\": \"mutation { merchantCreatePet(name:\\\"$name\\\", species:\\\"$species\\\", age:$age, pictureBase64:\\\"aGVsbG8=\\\", description:\\\"$description\\\") { id name species age status } }\"
    }"
  echo
}

create_pet "Luna" "DOG" 3 "Friendly golden retriever"
create_pet "Max" "DOG" 2 "Energetic beagle"
create_pet "Bella" "DOG" 4 "Loyal family dog"
create_pet "Charlie" "DOG" 1 "Playful puppy"
create_pet "Rocky" "DOG" 5 "Strong and calm dog"
create_pet "Daisy" "DOG" 2 "Sweet and active dog"
create_pet "Milo" "DOG" 3 "Curious brown dog"
create_pet "Cooper" "DOG" 6 "Smart shepherd mix"
create_pet "Sadie" "DOG" 4 "Happy outdoor dog"
create_pet "Toby" "DOG" 2 "Friendly small dog"

create_pet "Oliver" "CAT" 2 "Curious indoor cat"
create_pet "Leo" "CAT" 3 "Lazy orange cat"
create_pet "Lily" "CAT" 1 "Gentle white kitten"
create_pet "Nala" "CAT" 4 "Affectionate house cat"
create_pet "Simba" "CAT" 5 "Confident tabby cat"
create_pet "Chloe" "CAT" 2 "Playful gray cat"
create_pet "Lucy" "CAT" 3 "Calm and cuddly cat"
create_pet "Coco" "CAT" 1 "Cute young kitten"
create_pet "Zoe" "CAT" 6 "Independent senior cat"
create_pet "Jasper" "CAT" 4 "Elegant long-haired cat"

create_pet "Kiwi" "BIRD" 1 "Small bright parakeet"
create_pet "Sunny" "BIRD" 2 "Cheerful yellow canary"
create_pet "Sky" "BIRD" 1 "Blue bird with playful energy"
create_pet "Rio" "BIRD" 3 "Colorful tropical bird"
create_pet "Peanut" "BIRD" 2 "Friendly little bird"
