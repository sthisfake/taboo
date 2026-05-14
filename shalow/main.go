// main.go
package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const maxFolders = 10

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Auto Go Generator Started...")

	for {
		// STEP 1 - Cleanup old folders if needed
		cleanupOldFolders()

		// STEP 2 - Create folder
		folderName := "services_" + randomString(8)

		err := os.Mkdir(folderName, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating folder:", err)
			continue
		}

		fmt.Println("Created folder:", folderName)

		// STEP 3 - Random delay
		randomDelay()

		// STEP 4 - Create file
		fileName := "document_" + randomString(8) + ".go"
		filePath := filepath.Join(folderName, fileName)

		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Error creating file:", err)
			continue
		}

		fmt.Println("Created file:", filePath)

		// STEP 5 - Random delay
		randomDelay()

		// STEP 6 - Simulate long writing process
		writeDuration := randomDuration(1*time.Minute, 15*time.Minute)

		fmt.Println("Generating random Go code...")
		fmt.Println("This will take:", writeDuration)

		time.Sleep(writeDuration)

		code := generateLargeGoCode()

		_, err = file.WriteString(code)
		if err != nil {
			fmt.Println("Error writing code:", err)
			file.Close()
			continue
		}

		file.Close()

		fmt.Println("Finished writing Go code.")

		// STEP 7 - Random delay
		randomDelay()

		fmt.Println("Restarting loop...\n")
	}
}

func cleanupOldFolders() {
	entries, err := os.ReadDir(".")
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	type folderInfo struct {
		name    string
		modTime time.Time
	}

	var folders []folderInfo

	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "services_") {

			info, err := entry.Info()
			if err != nil {
				continue
			}

			folders = append(folders, folderInfo{
				name:    entry.Name(),
				modTime: info.ModTime(),
			})
		}
	}

	// If folder count is below threshold, do nothing
	if len(folders) < maxFolders {
		return
	}

	// Sort oldest first
	sort.Slice(folders, func(i, j int) bool {
		return folders[i].modTime.Before(folders[j].modTime)
	})

	// Delete oldest folders until under threshold
	deleteCount := len(folders) - maxFolders + 1

	for i := 0; i < deleteCount; i++ {
		folderToDelete := folders[i].name

		fmt.Println("Deleting old folder:", folderToDelete)

		err := os.RemoveAll(folderToDelete)
		if err != nil {
			fmt.Println("Failed to delete:", err)
		}
	}
}

func randomDelay() {
	delay := randomDuration(10*time.Second, 30*time.Second)

	fmt.Println("Sleeping for:", delay)

	time.Sleep(delay)
}

func randomDuration(min time.Duration, max time.Duration) time.Duration {
	diff := max - min
	return min + time.Duration(rand.Int63n(int64(diff)))
}

func randomString(length int) string {
	letters := "abcdefghijklmnopqrstuvwxyz"

	var builder strings.Builder

	for i := 0; i < length; i++ {
		builder.WriteByte(letters[rand.Intn(len(letters))])
	}

	return builder.String()
}

func generateLargeGoCode() string {
	var builder strings.Builder

	builder.WriteString("package main\n\n")

	builder.WriteString("import (\n")
	builder.WriteString(`	"fmt"` + "\n")
	builder.WriteString(`	"math"` + "\n")
	builder.WriteString(`	"time"` + "\n")
	builder.WriteString(")\n\n")

	builder.WriteString("type Worker struct {\n")
	builder.WriteString("	ID int\n")
	builder.WriteString("	Name string\n")
	builder.WriteString("}\n\n")

	builder.WriteString("func (w Worker) Print() {\n")
	builder.WriteString(`	fmt.Println("Worker:", w.ID, w.Name)` + "\n")
	builder.WriteString("}\n\n")

	// Generate many functions
	for i := 0; i < 25; i++ {

		builder.WriteString(fmt.Sprintf("func process%d(value int) int {\n", i))
		builder.WriteString("	result := value\n")

		for j := 0; j < 3; j++ {
			builder.WriteString(
				fmt.Sprintf(
					"	result += %d\n",
					rand.Intn(100),
				),
			)

			builder.WriteString(
				fmt.Sprintf(
					"	result -= %d\n",
					rand.Intn(50),
				),
			)

			builder.WriteString(
				"	result = int(math.Abs(float64(result)))\n",
			)
		}

		builder.WriteString("	return result\n")
		builder.WriteString("}\n\n")
	}

	builder.WriteString("func generateWorkers() []Worker {\n")
	builder.WriteString("	workers := []Worker{}\n")

	builder.WriteString("	for i := 0; i < 10; i++ {\n")
	builder.WriteString("		worker := Worker{\n")
	builder.WriteString("			ID: i,\n")

	builder.WriteString(
		`			Name: fmt.Sprintf("worker_%d", i),` + "\n",
	)

	builder.WriteString("		}\n")
	builder.WriteString("		workers = append(workers, worker)\n")
	builder.WriteString("	}\n")

	builder.WriteString("	return workers\n")
	builder.WriteString("}\n\n")

	builder.WriteString("func main() {\n")

	builder.WriteString(
		`	fmt.Println("Generated Go Program Running")` + "\n",
	)

	builder.WriteString("	workers := generateWorkers()\n")

	builder.WriteString("	for _, worker := range workers {\n")
	builder.WriteString("		worker.Print()\n")
	builder.WriteString("	}\n\n")

	builder.WriteString("	total := 0\n")

	for i := 0; i < 25; i++ {
		builder.WriteString(
			fmt.Sprintf(
				"	total += process%d(%d)\n",
				i,
				rand.Intn(500),
			),
		)
	}

	builder.WriteString(`	fmt.Println("Total:", total)` + "\n")

	builder.WriteString("	for i := 0; i < 5; i++ {\n")

	builder.WriteString(
		`		fmt.Println("Tick:", i)` + "\n",
	)

	builder.WriteString(
		"		time.Sleep(500 * time.Millisecond)\n",
	)

	builder.WriteString("	}\n")

	builder.WriteString(
		`	fmt.Println("Program Finished")` + "\n",
	)

	builder.WriteString("}\n")

	// Ensure roughly 80-110 lines
	extraLines := rand.Intn(20) + 10

	for i := 0; i < extraLines; i++ {
		builder.WriteString(
			fmt.Sprintf(
				"// extra generated line %d\n",
				i,
			),
		)
	}

	_ = math.MaxFloat64

	return builder.String()
}