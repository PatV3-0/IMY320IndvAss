package main

import (
    "math/rand"
    "syscall/js"
    "time"
)

func rollDice(sides int) int {
    return rand.Intn(sides) + 1
}

func rollD100() int {
    tens := rand.Intn(10) * 10
    ones := rand.Intn(10)
    result := tens + ones
    if result == 0 {
        result = 100
    }
    return result
}

func rollMultipleDice(num int, sides int) []int {
    if num < 1 {
        num = 1
    }

    results := make([]int, num)
    for i := 0; i < num; i++ {
        if sides == 100 {
            results[i] = rollD100()
        } else {
            results[i] = rollDice(sides)
        }
    }
    return results
}

func rollWrapped(this js.Value, args []js.Value) interface{} {
    if len(args) < 2 {
        return js.ValueOf([]interface{}{})
    }

    num := args[0].Int()
    sides := args[1].Int()

    validSides := map[int]bool{4: true, 6: true, 8: true, 10: true, 12: true, 20: true, 100: true}
    if !validSides[sides] {
        return js.ValueOf([]interface{}{})
    }

    results := rollMultipleDice(num, sides)

    // Convert Go []int to JS array
    jsArray := js.Global().Get("Array").New()
    for _, v := range results {
        jsArray.Call("push", v)
    }
    return jsArray
}

func main() {
    rand.Seed(time.Now().UnixNano())

    js.Global().Set("rollMultipleDice", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        return rollWrapped(this, args)
    }))

    select {}
}
