package main

import (
    "fmt"
    "strconv"
    "strings"
)


type Word struct {
 
    Name    string
    Message string
   
    Function func()
    Words   []float64
 
}


type Eval struct {

   Stack      Stack
   Dictionary []Word
   compiling  bool
   tmp        Word
  
   doOpen     int 
}

func NewEval() *Eval { 
    e := &Eval{}
   
   
e.Dictionary = []Word{
    {Name: "+",    Function: e.add},             
    {Name: "-",    Function: e.sub},             
    {Name: "*",    Function: e.mul},             
    {Name: "/",    Function: e.div},             
    {Name: ".",    Function: e.print},  
    {Name: ".s",   Function: e.printStack},
    {Name: ":",    Function: e.startDefinition},
    {Name: "cls",  Function: e.сlearStack},// чистит стек полностью
    {Name: "dup",  Function: e.dup},
    {Name: "show", Function: e.ShowDic},//выводит список слов
    {Name: "emit", Function: e.emit},
    {Name: "exit", Function: e.exit_prog},
   
    {Name: "do",   Function: e.do},
    {Name: "loop", Function: e.loop},
    
    

}
  return e
}


func (e *Eval) Eval(args []string) {

  /*var error_str int = 0 */ 
    
  for _, tok := range args {
      
    tok = strings.TrimSpace(tok)
    
    if tok == "" {
        continue
    }

      if tok == ":" {
              e.startDefinition()  // e.compiling == true
              continue
         }

    if e.compiling {
    
        if e.tmp.Name == "" {
           
            idx := e.findWord(tok)
           
            if idx != -1 {
                fmt.Printf("word %s already defined\n", tok)
                return
            }

            // сохраняем имя после символа :
            e.tmp.Name = tok
            continue
        }
//////////////////////////////// //////////Проверяет на сообщение и записывает его
          // Возможно добавить или изменить на ."
            if tok == ".'"  {
                             
                s := strings.Join(args, " ")
              
                i := strings.Index(s, ".'")
//                     fmt.Println("index: ",i) Точка начала сообщения
                            
                ib := strings.Index(s, ";")
//                     fmt.Println("index: ",ib) Конечная точка начала сообщения
                
                a := s[(i+2):(ib)]
                   
                e.tmp.Message = a            
                /*error_str = 1   */  
                 
                fmt.Println("Add Message :", e.tmp.Message)
                continue              
            }
            
            
        

        if tok == ";"/* && error_str != 1*/{ // конец создания переменной
            e.Dictionary = append(e.Dictionary, e.tmp)
            e.tmp.Name = ""
            e.tmp.Words = []float64{}
            e.compiling = false
                     
            e.tmp.Message = ""
           
           
            continue

        } // компиляция

     
        idx := e.findWord(tok) // вовзращает найденный индекс слова в словаре
     
        
       if idx >= 0 {
         
           e.tmp.Words = append(e.tmp.Words, float64(idx))  //Добавляет индекс в массив для слов
           
           if tok == "do" {
            
               e.doOpen = len(e.tmp.Words) - 1
           }
           
           if tok == "loop" {
               
             e.tmp.Words = append(e.tmp.Words, -2)
             e.tmp.Words = append(e.tmp.Words, float64(e.doOpen))
               
           }           
           
           
       } else { // если нет то предполагаем, что число
      
           
            e.tmp.Words = append(e.tmp.Words, -1)
            
            if e.tmp.Message != "" {
                 fmt.Printf("Message : %s\n",e.tmp.Message)
              continue   
             }
            
            
            val, err := strconv.ParseFloat(tok, 64)
            if err != nil {
                fmt.Printf("OOO %s: %s\n", tok, err.Error())
                      
            }            
            
            e.tmp.Words = append(e.tmp.Words, val)
             
        }

        continue
    }

//    Обработали как элемент словаря?
    handled := false
    for index, word := range e.Dictionary {
        if tok == word.Name {
            e.evalWord(index)//Если слово то вычисляем
            
            handled = true
        }
    }


    if !handled { // Если не слово то определяем как число
        i, err := strconv.ParseFloat(tok, 64)
       
        if err != nil {
            fmt.Println("Нет в словаре, добавьте в словарь\n")
            fmt.Printf("Error: %s: %s\n", tok, err.Error())// если нет в словаре
            return
        }

        e.Stack.Push(i)
    }
}

}

// evalWord Вычисляет слово по индексу из словаря
func (e *Eval) evalWord(index int) {

   word := e.Dictionary[index]
     
   // Выводи индекс слова и сообщение в нем содержащиеся  
   if word.Function == nil  {
            
       if word.Message != "" {
       
       fmt.Printf("Message : %s\n", word.Message)
       }
          
   }// Конец вывода сообщения
   
   if word.Function != nil  {
       
       word.Function()
         
   } 
   /////////////////////////////
   
   addNum := false
   
   jump   := false
   
   inst   := 0 // Хранилище индексов
   
   for inst < len(word.Words) {
       
       opcode := word.Words[inst] // Текущая операция
       
       if addNum { // Если число
            
            e.Stack.Push(opcode)
            addNum = false
     
           
       } else if jump {
           
            cur := e.Stack.Pop() 
            max := e.Stack.Pop() 
         
            if max > cur {
             
               e.Stack.Push(max)
               e.Stack.Push(cur)  
             
               inst = int(opcode) 
               
               inst--
            
            }
                 
         jump = false
           
         } else {

            // есди видим -1 ,то добавляем число
            if opcode == -1 {
               
                addNum = true
                            
            } else if opcode == -2 {
                
                jump = true  
                
              } else {

                  // иначе как обычно
                  e.evalWord(int(opcode))
                }
        } 
        
        inst++
       
   }
  
}

func (e *Eval) findWord(name string) int {
   for index, word := range e.Dictionary {
       if name == word.Name {
           return index
       }
   }
   return -1
}
