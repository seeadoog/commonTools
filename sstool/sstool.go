package main

import (
	"flag"
	"io"
	"os"
)

type Command struct {
	Name string
	Flag *flag.FlagSet
	Run func()
}

type App struct {
	Args []string
	Commands []Command
	Flag *flag.FlagSet
	Runa func()
}
func NewApp()*App{
	return &App{
		Args:     os.Args,
		Commands: nil,
		Flag:     flag.NewFlagSet(os.Args[0],flag.ExitOnError),
	}
}
func (a *App)Run()error{
	if len(a.Args) <=1{
		a.Runa()
	}
	cmd:=a.Args[1]
	for _, v := range a.Commands {
		if v.Name == cmd{
			if err:=v.Flag.Parse(a.Args[2:]);err != nil{
				return err
			}
			v.Run()
		}
	}
	if err:=a.Flag.Parse(a.Args[1:]);err != nil{
		return err
	}

	a.Runa()
	return nil
}

type Flag struct {
	Name string
}




func read(in string)(io.Reader,error){
	if in== ""{
		return os.Stdin,nil
	}

	return os.Open(in)
}

func main(){

	in,err:=read(os.Args[1])


}
