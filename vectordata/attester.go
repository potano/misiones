// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import "potano.misiones/sexp"


const (
	weightedAttestationGroup = iota + 1
	singleValuedAttestationGroup
)

type attestationGroup struct {
	name string
	groupType int
	groupID int
	sumWeights int
	millsPerStep int
}

type attestationDef struct {
	groupNum int
	weight int
}

type attester struct {
	doc *VectorData
	groups []attestationGroup
	allowedAttestations map[string]attestationDef
}


func newAttester(doc *VectorData) *attester {
	return &attester{
		doc: doc,
		allowedAttestations: map[string]attestationDef{},
	}
}

func (att *attester) checkConfiguration() error {
	return nil
}

func (att *attester) resolveStyle(attestation *mapAttestationType, style *mapStyleType) error {
	var styleX int
	var err error
	if style != nil {
		styleX, err = att.doc.styler.baseStyleIndex(style)
		if err != nil {
			return err
		}
	}
	groupUse := make([]int, len(att.groups))
	oneTimeUsers := make([]string, len(att.groups))
	for _, name := range attestation.attestations {
		if attDef, exists := att.allowedAttestations[name]; !exists {
			return attestation.Error("unknown attestation '%s'")
		} else {
			groupNum := attDef.groupNum
			groupInfo := att.groups[groupNum]
			if groupInfo.groupType == weightedAttestationGroup {
				groupUse[groupNum] += attDef.weight
			} else if groupInfo.groupType == singleValuedAttestationGroup {
				if len(oneTimeUsers[groupNum]) > 0 {
					return attestation.Error(
						"multiple %s attestations; cannot use %s with %s",
						groupInfo.name, name, oneTimeUsers[groupNum])
				}
				// offset index by 1 so style index=0 indicates 'no added properties'
				groupUse[groupNum] = attDef.weight + 1
			}
			oneTimeUsers[groupNum] = name
		}
	}
	for groupNum, groupInfo := range att.groups {
		if groupInfo.groupType == weightedAttestationGroup {
			if len(oneTimeUsers[groupNum]) > 0 {
				groupMills := groupUse[groupNum] * 1000
				groupUse[groupNum] = (groupMills / groupInfo.millsPerStep) + 1
			}
		}
	}
	attestation.resolvedStyleIndex = att.doc.styler.findAttestationStyle(styleX, groupUse)
	return nil
}





type mapAttestationType struct {
	mapItemCore
	attestations []string
	resolvedStyleIndex int
}

func newMapAttestation(doc *VectorData, parent mapItemType, listType, listName string,
		source sexp.ValueSource) (mapItemType, error) {
	ms := &mapAttestationType{}
	ms.itemType = mitAttestation
	ms.source = source
	return ms, nil
}

func (ma *mapAttestationType) addScalars(targetName string, scalars []sexp.LispScalar) error {
	ma.attestations = make([]string, len(scalars))
	seen := make(map[string]bool, len(scalars))
	for i, scalar := range scalars {
		attestation := scalar.String()
		if seen[attestation] {
			return scalar.Error("duplicate attestation")
		}
		ma.attestations[i] = attestation
	}
	return nil
}

