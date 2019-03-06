 #!/bin/bash

fGetSizeByHuman()
 {
   sizeStr=""
   sum=$1
   if [ $sum -ge $GB ]; then
      gb=$(awk 'BEGIN{printf "%.2f\n", '$sum'/'$GB'}')
      sizeStr="$gb GB"
   elif [ $sum -ge $MB ]; then
      mb=$(awk 'BEGIN{printf "%.2f\n", '$sum'/'$MB'}')
      sizeStr="$mb MB"
   elif [ $sum -ge $KB ]; then
      kb=$(awk 'BEGIN{printf "%.2f\n", '$sum'/'$KB'}')
      sizeStr="$kb KB"
   else
      sizeStr="$sum B"
   fi
   echo $sizeStr
 }

 index=0
 todir="/todir-"
 echo "" > /fail.txt
 echo "" > /result.txt

 (( GB = 1024 * 1024 * 1024 ))
 (( MB = 1024 * 1024 ))
 (( KB = 1024 ))

 for dir in $*
 do
   scpto=$todir$index
   (( index = index +1 ))
   {
      echo "start scp $dir $scpto"
      sleep 1
      rsync -ravP $dir $scpto 1>/dev/null 2>/dev/null
      if [ "$?" == "0" ]; then
         echo "stop scp $dir $scpto success"
      else
         echo "true" > /fail.txt
         echo "stop scp $dir $scpto fail"
      fi
      echo "stop scp $dir $scpto" >> /result.txt
   }&
 done

 {
    lastsum=0
    while true
    do
       sum=0
       index=0
       for dir in $*
       do
          dirname=${dir##*/}

         scpto=$todir$index/$dirname
          (( index = index + 1 ))
          size=`du -b --max-depth=0 $scpto 2>/dev/null | awk '{print $1}'`
          (( sum = sum + size ))
       done

       (( diff= sum - lastsum ))
       lastsum=$sum

       echo "Has moved $(fGetSizeByHuman $sum) ($(fGetSizeByHuman $diff)/s)"
       if [ "$(cat /result.txt | grep stop | wc -l)" == "$index" ]; then
          exit 0
       fi
       sleep 1
    done
 }&

 wait

 if [ "$(cat /fail.txt)" == "true" ]; then
    echo "scp fail"
    exit 1
 else
    echo "scp success"
    exit 0
 fi
