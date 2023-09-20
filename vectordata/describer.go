// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import (
	"fmt"
	"strings"
)


func (vd *VectorData) DescribeNodes(indent string) string {
	return describeMapItem("", indent, vd.layersRoot)
}


func summarizeMapItem(pad string, item mapItemType) string {
	return fmt.Sprintf("%s%s '%s'", pad, typeMapToName[item.ItemType()], item.Name())
}

func describeMapItem(pad, indent string, item mapItemType) string {
	lines := []string{
		fmt.Sprintf("%s→%s '%s' @ %s", pad, typeMapToName[item.ItemType()], item.Name(),
			item.Source().String())}
	var children []mapItemType
	padpad := pad + "    "
	switch item := item.(type) {
	case *mapLayersType:
		children = make([]mapItemType, len(item.layers))
		for i, layer := range item.layers {
			children[i] = layer
		}
	case *mapLayerType:
		lines = append(lines, pad + "    menuitem: '" + item.menuitem + "'")
		children = item.features
	case *mapFeatureType:
		describePopup(&lines, padpad, item.popup)
		describeStyle(&lines, padpad, item.style)
		describeAttestation(&lines, padpad, item.attestation)
		children = item.features
	case *mapMarkerType:
		describePopup(&lines, padpad, item.popup)
		if len(item.html) > 0 {
			lines = append(lines, padpad + "html: '" + stringUpTo(25, item.html) + "'")
		}
		describeLocation(&lines, padpad, item.location)
	case *map_referenceAggregateType:
		lines = append(lines, padpad + "parent: " + item.parentName)
		line := padpad + "target names:"
		for _, targ := range item.names {
			name := targ.String()
			if len(line) + len(name) > 78 {
				lines = append(lines, line)
				line = padpad + "             "
			}
			line += " " + name
		}
		lines = append(lines, line)
		children = item.targets
	case *map_locationType:
		describePopup(&lines, padpad, item.popup)
		describeStyle(&lines, padpad, item.style)
		describeAttestation(&lines, padpad, item.attestation)
		if item.ItemType() == mitCircle {
			units := "meters"
			if item.radiusType == mitPixels {
				units = "pixels"
			}
			lines = append(lines, fmt.Sprintf("%sradius: %d %s", padpad, item.radius,
				units))
		}
		describeLocation(&lines, padpad, item.location)
	case *mapRouteType:
		describePopup(&lines, padpad, item.popup)
		describeStyle(&lines, padpad, item.style)
		describeAttestation(&lines, padpad, item.attestation)
		children = item.segments
	case *mapSegmentType:
		describePopup(&lines, padpad, item.popup)
		describeStyle(&lines, padpad, item.style)
		describeAttestation(&lines, padpad, item.attestation)
		children = item.paths
	default:
		tp := item.ItemType()
		panic(fmt.Sprintf("describing unexpected type %s (%d)", typeMapToName[tp], tp))
	}
	for _, child := range children {
		lines = append(lines, describeMapItem(pad + indent, indent, child))
	}
	return strings.Join(lines, "\n")
}


func describePopup(lines *[]string, pad string, popup *mapPopupType) {
	if popup == nil {
		return
	}
	*lines = append(*lines, fmt.Sprintf("%spopup text: '%s'", pad, popup.text))
}

func describeStyle(lines *[]string, pad string, style *mapStyleType) {
	if style == nil {
		return
	}
	*lines = append(*lines, pad + "style: " + style.name)
}

func describeAttestation(lines *[]string, pad string, attestation *mapAttestationType) {
	if attestation == nil {
		return
	}
	*lines = append(*lines, pad + "attestation: " + strings.Join(attestation.attestations, " "))
}

func describeLocation(lines *[]string, pad string, location locationPairs) {
	label := "location: "
	for i := 0; i < len(location); i += 2 {
		*lines = append(*lines, fmt.Sprintf("%s%-10s%.6f  %.6f", pad, label, location[i],
			location[i+1]))
		label = ""
	}
}

func stringUpTo(maxlen int, str string) string {
	runes := []rune(str)
	if len(runes) <= maxlen {
		return str
	}
	return string(runes[0:maxlen-3]) + "..."
}

