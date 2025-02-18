package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const version = "0.1.2"

type episode struct {
	fileName string

	number int
}

func main() {
	showVersion := flag.Bool("version", false, "Show version")
	series := flag.String("series", "", "Series name")
	season := flag.Int("season", 0, "Season number")
	forceFolder := flag.Bool("force_folder", false, "Proceed even if folder already exists")
	replaceTitle := flag.Bool("replace_title", false, "Replate the metadata title with file name")
	flag.Parse()

	if *showVersion {
		fmt.Printf("jellyfy version %s\n", version)
		return
	}

	folder := flag.Arg(0)

	if folder == "" || *series == "" || *season == 0 {
		fmt.Println("Usage: jellyfy -series=<series> -season=<season> <folder>")
		return
	}

	targetFolderName := fmt.Sprintf("Season %02d", *season)

	newFolderPath, err := renameFolder(folder, targetFolderName, *forceFolder)
	if err != nil {
		fmt.Println("Error renaming folder:", err)
		return
	}

	episodes, err := listEpisodeFiles(newFolderPath)
	if err != nil {
		fmt.Println("Error listing episodes:", err)
		return
	}

	sort.Slice(episodes, func(i, j int) bool {
		return episodes[i].number < episodes[j].number
	})

	for _, episode := range episodes {
		ext := strings.ToLower(filepath.Ext(episode.fileName))
		targetName := fmt.Sprintf("%s S%02dE%02d", *series, *season, episode.number)
		targetFileName := targetName + ext

		srcPath := filepath.Join(newFolderPath, episode.fileName)
		newPath := filepath.Join(newFolderPath, targetFileName)
		if targetFileName != episode.fileName {
			if err := os.Rename(srcPath, newPath); err != nil {
				fmt.Println("Error renaming episode:", err)
			}
			fmt.Printf("%s -> %s\n", episode.fileName, targetFileName)
		}

		if ext == ".mkv" && *replaceTitle {
			if err := modifyMkvTitle(newPath, targetName); err != nil {
				fmt.Println("Error modifying title:", err)
			}
			fmt.Printf("Title modified in %s\n", targetFileName)
		}
	}
}

func renameFolder(folderPath string, newName string, force bool) (string, error) {
	parentDir := filepath.Dir(folderPath)
	newPath := filepath.Join(parentDir, newName)

	// Check if target already exists
	if _, err := os.Stat(newPath); err == nil {
		if !force {
			return "", fmt.Errorf("target folder already exists: '%s'", newName)
		}
		return newPath, nil
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("error checking target folder: %w", err)
	}

	if err := os.Rename(folderPath, newPath); err != nil {
		return "", fmt.Errorf("error renaming folder: %w", err)
	}
	fmt.Println("Folder renamed successfully")

	return newPath, nil
}

func listEpisodeFiles(folderPath string) ([]episode, error) {
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %w", err)
	}

	// Match patterns like S01E01, s1e1, S1E01, etc and capture episode number
	episodePattern := regexp.MustCompile(`(?i)s\d+e(\d+)`)
	var episodes []episode

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if matches := episodePattern.FindStringSubmatch(name); matches != nil {
			episodeNum, _ := strconv.Atoi(matches[1])
			episodes = append(episodes, episode{
				fileName: name,

				number: episodeNum,
			})
		}
	}

	return episodes, nil
}

func modifyMkvTitle(filepath string, newTitle string) error {
	// To remove title completely, use --delete title
	// To set new title, use --set title="new title"
	args := []string{"--edit", "info", "--set", "title=" + newTitle}
	if newTitle == "" {
		args = []string{"--edit", "info", "--delete", "title"}
	}

	cmd := exec.Command("mkvpropedit", append([]string{filepath}, args...)...)
	return cmd.Run()
}
