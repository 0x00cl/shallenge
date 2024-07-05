// shallenge - find the lowest hash possible for a given string format
// Copyright (C) 2024  Tomás Gutiérrez L. (0x00)

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

const (
	version      = "0.1.0"
	license_text = `Copyright (C) 2024 Tomás Gutiérrez L. (0x00)
License GPLv3: GNU GPL version 3 <https://gnu.org/licenses/gpl.html>.
This is free software: you are free to change and redistribute it.
This program comes with ABSOLUTELY NO WARRANTY.
`
	string_length = 18
	letter_bytes  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+/"
)

var initial_string = "0x00cl/"
var suffix_string string
var worker_count int

func generate_random_string() string {
	b := make([]byte, string_length)
	for i := range b {
		b[i] = letter_bytes[rand.Int63()%int64(64)]
	}
	return string(b)
}

func worker(wg *sync.WaitGroup, result_chan chan<- string, lowest_hash *string, lowest_hash_mutex *sync.Mutex) {
	defer wg.Done()

	for {
		random_string := generate_random_string()
		string_to_hash := initial_string + random_string
		hash := sha256.Sum256([]byte(string_to_hash))
		hash_string := hex.EncodeToString(hash[:])
		lowest_hash_mutex.Lock()
		if hash_string < *lowest_hash {
			*lowest_hash = hash_string
			t := time.Now()
			result_chan <- "[" + t.Format("2006-01-02T15:04:05 -070000") + "] " + hash_string + " (" + string_to_hash + ")"
		}
		lowest_hash_mutex.Unlock()
	}
}

func init() {
	flag.StringVar(&suffix_string, "p", "i5+8250U/Hello+HN/", "String to add as suffix to the random generated word")
	flag.IntVar(&worker_count, "n", 8, "Number of workers to spawn")
}

func main() {
	version_flag := flag.Bool("v", false, "Print the version number.")
	help_flag := flag.Bool("h", false, "Print this help message.")

	flag.Parse()

	if *version_flag {
		fmt.Printf("shallenge %s\n%s", version, license_text)
		os.Exit(0)
	}

	if *help_flag {
		flag.PrintDefaults()
		os.Exit(0)
	}

	var wg sync.WaitGroup
	result_chan := make(chan string)
	lowest_hash := "00000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
	var lowest_hash_mutex sync.Mutex

	initial_string = initial_string + suffix_string
	for i := 0; i < worker_count; i++ {
		wg.Add(1)
		go worker(&wg, result_chan, &lowest_hash, &lowest_hash_mutex)
	}

	go func() {
		wg.Wait()
		close(result_chan)
	}()

	for result := range result_chan {
		fmt.Println(result)
	}
}
