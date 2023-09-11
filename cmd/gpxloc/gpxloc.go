package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/urfave/cli"
)

type Config struct {
	latitude  float64
	longitude float64
	radius    int64
	bbox      string
	filename  string
	// args are the positional (non-flag) command-line arguments.
	args []string
}

func main() {
	conf := Config{}

	app := cli.NewApp()
	app.Name = "gpxloc"
	app.Usage = "GPX files Location Finder\n\n   Find gpx files according to the location where they were recorded"
	app.Version = "0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "bbox",
			Usage:       "bounding box (\"lon1,lat1,lon2,lat2\") ",
			Required:    false,
			Destination: &conf.bbox,
		},

		cli.Float64Flag{
			Name:        "latitude, lat",
			Usage:       "latitude (WGS84 [-90,+90])",
			Required:    false,
			Destination: &conf.latitude,
		},
		cli.Float64Flag{
			Name:        "longitude, lon",
			Usage:       "longitude (WGS84 [-180,+180])",
			Required:    false,
			Destination: &conf.longitude,
		},
		cli.Int64Flag{
			Name:        "radius",
			Usage:       "radius (in meters)",
			Required:    false,
			Destination: &conf.radius,
		},
	}

	app.Action = func(c *cli.Context) error {
		conf.args = c.Args()
		err := validateConfig(&conf)
		if err != nil {
			return err
		}
		doWork(&conf)
		return nil
	}

	cli.AppHelpTemplate = fmt.Sprintf(`%s

WEBSITE: https://github.com/JVillafruela/GPXloc

EXAMPLES:
   gpxloc --bbox="5.68678,45.08596,5.68979,45.08778" E:\OSM\gps\2022\2022-04-16 E:\OSM\gps\2022\2022-04-22

   gpxloc --latitude=45.087 --longitude=5.688 --radius=20  E:\OSM\gps\2023

	`, cli.AppHelpTemplate)

	app.Setup()
	// do not display command in help (see cli issue #523)
	app.Commands = []cli.Command{}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// check if the configuration is valid
func validateConfig(config *Config) error {

	bboxSet := config.bbox != ""
	pointSet := config.latitude != 0 && config.longitude != 0 && config.radius != 0

	if !bboxSet && !pointSet {
		return errors.New("indicate either the bbox option or the lat,lon,radius or the radius options")
	}

	if bboxSet && pointSet {
		return errors.New("too many options. Indicate either the bbox option or the lat,lon,radius options or the radius options")
	}

	if bboxSet {
		_, err := BboxBound(config.bbox)
		if err != nil {
			return errors.New("invalid value for bounding box option. Should be \"lon1,lat1,lon2,lat2\" ")
		}
	}

	if pointSet {
		err := validateLatitude(config.latitude)
		if err != nil {
			return err
		}
		err = validateLongitude(config.longitude)
		if err != nil {
			return err
		}
		if config.radius <= 0 {
			return errors.New("invalid radius")
		}
	}

	if len(config.args) == 0 {
		return errors.New("missing argument")
	}

	for _, v := range config.args {
		path := cleanPath(v)
		if !dirExists(path) {
			return errors.New("Directory does not exist : " + path)
		}
	}

	return nil
}

// process the files
func doWork(conf *Config) error {
	//fmt.Printf("config = %+v\n", *config)
	var lf LocationFinder
	var err error
	if conf.bbox != "" {
		lf, err = NewBboxLocationFinder(conf.bbox)
		if err != nil {
			return err
		}
	}

	if conf.latitude != 0 && conf.longitude != 0 && conf.radius != 0 {
		lf, err = NewCircleLocationFinder(conf.latitude, conf.longitude, conf.radius)
		if err != nil {
			return err
		}
	}

	for _, dir := range conf.args {
		err = walkGPXFilesInDir(dir, lf)
	}

	return err
}

// traverse a directory looking for geotagged files whose position is inside the bounding box
func walkGPXFilesInDir(dir string, lf LocationFinder) error {
	return filepath.WalkDir(dir, func(path string, file fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := file.Info()
		if err != nil {
			return err
		}

		if !file.IsDir() && hasGPXExtension(info.Name()) {
			fullname := cleanPath(path)

			gpxfile, err := ReadGPXfile(fullname)
			if err != nil {
				fmt.Println("File name: ", fullname, "error ", err)
				return nil
			}
			//fmt.Println("DDD file name: ", fullname, gpxfile.bounds)

			if lf.PossibleMatch(gpxfile.bounds) {
				for _, track := range gpxfile.gpx.Tracks {
					for _, segment := range track.Segments {
						for _, point := range segment.Points {
							if lf.Match(point.Latitude, point.Longitude) {
								fmt.Println(fullname)
								return nil
							}
						}
					}
				}

			}
		}
		return nil
	})
}

// check if the file has a jpeg extension
func hasGPXExtension(filename string) bool {
	extensions := []string{"gpx"} // sorted
	return hasExtension(filename, extensions)
}

func hasExtension(filename string, extensions []string) bool {
	if len(extensions) == 0 {
		return false
	}

	ext := strings.ToLower(strings.Trim(filepath.Ext(filename), "."))
	if ext == "" {
		return false
	}
	i := sort.SearchStrings(extensions, ext)
	return i < len(extensions) && extensions[i] == ext
}

// on Windows get rid of the trailing " added by autocompletion \"
func cleanPath(path string) string {
	if runtime.GOOS == "windows" {
		path = strings.TrimSuffix(path, `"`)
	}
	return filepath.Clean(path)
}

// check if filename is a directory
func dirExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dirExists "+filename, err.Error())
		return false
	}
	if os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}

// check if filename is a regular file
func fileExists(filename string) bool {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}
