// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package main

import (
	"os"
	"fmt"
	"flag"
	"path/filepath"

	"potano.misiones/sexp"
	"potano.misiones/great"
	"potano.misiones/vectordata"
)

func main() {
	sourceDir := "."
	var generateFile, measureName string
	var upToDistance float64
	var asMiles bool

	flag.StringVar(&sourceDir, "d", ".", "directory containing .sexp files")
	flag.StringVar(&generateFile, "g", "", "name of target Javascript file")
	flag.StringVar(&measureName, "m", "", "name of path or route to measure")
	flag.Float64Var(&upToDistance, "u", 0.0,
		"measure path only up to distance; report coordinates")
	flag.BoolVar(&asMiles, "miles", false, "measure distances in miles, not meters")
	flag.Parse()

	if !isDir(sourceDir) {
		fatal("source directory %s does not exist", sourceDir)
	}

	vd := vectordata.NewVectorData()
	vdReader, err := vectordata.NewVectorDataReader(vd)
	if err != nil {
		fatal(err.Error())
	}
	names, err := filepath.Glob(sourceDir + "/*.sexp")
	if err != nil {
		fatal(err.Error())
	}
	for _, filename := range names {
		fh, err := os.Open(filename)
		if err != nil {
			fatal(err.Error())
		}
		sourceList, err := sexp.Parse(filename, fh)
		if err != nil {
			fatal(err.Error())
		}
		err = vdReader.ConsumeList(sourceList)
		if err != nil {
			fatal(err.Error())
		}
	}
	err = vd.ResolveReferences()
	if err != nil {
		fatal(err.Error())
	}

	if len(generateFile) > 0 {
		blob, err := vd.GenerateJs()
		if err != nil {
			fatal(err.Error())
		}
		var outfile *os.File
		if generateFile == "-" {
			outfile = os.Stdout
		} else {
			outfile, err = os.Create(generateFile)
			if err != nil {
				fatal(err.Error())
			}
		}
		_, err = outfile.Write([]byte(blob))
		if err != nil {
			fatal(err.Error())
		}
		outfile.Close()
	}

	if len(measureName) > 0 {
		if upToDistance < 0 {
			fatal("argument to the -u switch must not be negative")
		}
		if upToDistance == 0 {
			distance, err := vd.MeasurePath(measureName)
			if err != nil {
				fatal(err.Error())
			}
			fmt.Printf("%0.1f meters (%0.2f miles)\n", distance,
				distance / great.METERS_PER_MILE)
		} else {
			if asMiles {
				upToDistance *= great.METERS_PER_MILE
			}
			lat, long, distance, pathName, index, err := vd.MeasurePathUpTo(
				measureName, upToDistance)
			if err != nil {
				fatal(err.Error())
			}
			fmt.Printf("Distance to latitude %.6f, longitude %.6f: %.1f meters " +
				"(%.1f miles)\n at point %d along path %s\n", lat, long, distance,
				distance / great.METERS_PER_MILE, index, pathName)
		}
	}
}



func isDir(filename string) bool {
	fi, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return fi.IsDir()
}


func fatal(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg + "\n", args...)
	os.Exit(1)
}

