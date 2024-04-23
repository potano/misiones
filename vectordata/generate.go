// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import (
	"strings"
	"strconv"
)

func (vd *VectorData) GenerateJs() (string, error) {
	if !vd.styler.styleCheckRun() {
		err := vd.CheckInStylesAndAttestations()
		if err != nil {
			return "", err
		}
	}
	var blobs []string
	if vd.styler != nil {
		blobs = append(blobs, vd.styler.generateJs())
	}
	for _, name := range vd.inDependencyOrder {
		obj := vd.mapItems[name]
		if len(obj.Referrers()) > 1 && obj.ItemType() != mitPoint {
			blobs = append(blobs, "var " + obj.Name() + "=" + obj.generateJs())
		}
	}
	blobs = append(blobs, "allVectors=" + vd.layersRoot.generateJs())
	return "(function() {" + strings.Join(blobs, "\n") + "})();", nil
}


func (ml *mapLayersType) generateJs() string {
	asFeatures := make([]mapItemType, len(ml.layers))
	for i, v := range ml.layers {
		asFeatures[i] = v
	}
	return featurizer(asFeatures).generateJs()
}


func (ml *mapLayerType) generateJs() string {
	return generateJsObject("menuitem", ml.menuitem, "features", featurizer(ml.features))
}


func (mf *mapFeatureType) generateJs() string {
	return generateJsObject(
		"t", "feature",
		"popup", mf.popup.nonEmptyString(),
		"style", attestationOrStyle(mf.attestation, mf.style),
		"features", featurizer(mf.features))
}


func (mp *mapPopupType) nonEmptyString() nonEmptyString {
	var text string
	if mp != nil {
		text = mp.text
	}
	return nonEmptyString(text)
}


func (ml *map_locationType) generateJs() string {
	if ml.itemType == mitCircle {
		var asPixels bool
		if ml.radiusType == mitPixels {
			asPixels = true
		}
		return generateJsObject(
			"t", "circle",
			"popup", ml.popup.nonEmptyString(),
			"style", attestationOrStyle(ml.attestation, ml.style),
			"asPixels", asPixels,
			"radius", ml.radius,
			"coords", ml.location)
	}
	return generateJsObject(
		"t", typeMapToName[ml.itemType],
		"popup", ml.popup.nonEmptyString(),
		"html", nonEmptyString(ml.html),
		"style", attestationOrStyle(ml.attestation, ml.style),
		"coords", ml.location)
}


func (mr *mapRouteType) generateJs() string {
	return generateJsObject(
		"t", "route",
		"popup", mr.popup.nonEmptyString(),
		"style", attestationOrStyle(mr.attestation, mr.style),
		"features", featurizer(mr.segments))
}


func (ms *mapSegmentType) generateJs() string {
	return generateJsObject(
		"t", "segment",
		"popup", ms.popup.nonEmptyString(),
		"style", attestationOrStyle(ms.attestation, ms.style),
		"paths", featurizer(ms.paths))
}


func (mr *map_referenceAggregateType) generateJs() string {
	if mr.itemType == mitRouteSegments {
		return generateJsFeaturesList(mr.targets[0].(*mapRouteType).segments)
	}
	return generateJsFeaturesList(mr.targets)
}



type nonEmptyString string
type nonEmptyCode string

type featurizer []mapItemType

func (mis featurizer) generateJs() string {
	return "[" + generateJsFeaturesList(mis) + "]"
}



func generateJsFeaturesList(features []mapItemType) string {
	blobs := make([]string, len(features))
	var i int
	for _, feature := range features {
		if feature.ItemType() == mitPoint {
			continue;
		}
		if len(feature.Referrers()) > 1 {
			blobs[i] = feature.Name()
		} else {
			blobs[i] = feature.generateJs()
		}
		i++
	}
	return strings.Join(blobs[:i], ",")
}

func generateJsObject(args ...any) string {
	entries := []string{}
	for i := 0; i < len(args) - 1; i += 2 {
		key := args[i]
		val := args[i+1]
		var str string
		switch v := val.(type) {
		case int:
			str = strconv.Itoa(v)
		case float64:
			str = strconv.FormatFloat(v, 'f', 6, 64)
		case locAngleType:
			str = v.String()
		case string:
			str = strconv.Quote(v)
		case nonEmptyString:
			if len(v) == 0 {
				continue
			}
			str = strconv.Quote(string(v))
		case nonEmptyCode:
			if len(v) == 0 {
				continue
			}
			str = string(v)
		case bool:
			str = strconv.FormatBool(v)
		case mapItemType:
			if val == nil {
				continue
			}
			str = v.generateJs()
		case locationPairs:
			str = v.generateJs()
		case featurizer:
			str = v.generateJs()
		default:
			continue
		}
		entries = append(entries, strconv.Quote(key.(string)) + ":" + str)
	}
	return "{" + strings.Join(entries, ",") + "}"
}


func attestationOrStyle(attestation *mapAttestationType, style *mapStyleType) nonEmptyCode {
	var text string
	resolvedStyleIndex := -1
	if attestation != nil {
		resolvedStyleIndex = attestation.resolvedStyleIndex
	} else if style != nil {
		resolvedStyleIndex = style.resolvedStyleIndex
	}
	if resolvedStyleIndex > 0 {
		text = formStyleName(resolvedStyleIndex)
	}
	return nonEmptyCode(text)
}

