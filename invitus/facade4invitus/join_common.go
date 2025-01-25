package facade4invitus

//func stringHash(s string) (hash int32) {
//	if len(s) == 0 {
//		return
//	}
//	for _, char := range s {
//		hash = ((hash << 5) - hash) + char
//	}
//	return
//}

//func getPin(spaceID, role, uid string) (pin int32) {
//	if pin = stringHash(spaceID + "-" + role + "-" + uid); pin < 0 {
//		pin = -pin
//	}
//	return
//}

//func verifyPin(spaceID string, pin int32, members []*briefs4memberus.MemberBrief) (role string, inviter *briefs4memberus.MemberBrief) {
//	for _, m := range members {
//		if m.UserID != "" {
//			if pin == getPin(spaceID, briefs4memberus.SpaceMemberRoleContributor, m.UserID) {
//				return briefs4memberus.SpaceMemberRoleContributor, m
//			}
//			if pin == getPin(spaceID, briefs4memberus.SpaceMemberRoleSpectator, m.UserID) {
//				return briefs4memberus.SpaceMemberRoleSpectator, m
//			}
//		}
//	}
//	return "", nil
//}
