package chunker

import (
	"log"
	"github.com/fsnotify/fsnotify"
)

//usinbg fsnotify for file chnage and trigger synchronization
//when file is modified
func WatchFile(filename string, changefilechannel chan bool){
   watcher,err:=fsnotify.NewWatcher()
   if err!=nil{
      log.Fatal(err)
   }
   //close watcher when the fxn exits
   defer watcher.Close()
   err=watcher.Add(filename)
   if err!=nil{
      log.Fatal(err)
   }
   for{
      select{
      case event,ok:=<-watcher.Events:
         if !ok{
            return
         }
         //check events corresponds to write operation on file
         if event.Op&fsnotify.Write==fsnotify.Write{
            log.Println("Modified file",event.Name)
            //send signal to chnage channel indicate file modification
            changefilechannel<-true
         }
      case err,ok:=<-watcher.Errors:
         if !ok{
            return 
            //exit fxn if channel is closed
         }
         log.Println(err)
      }
   }
}