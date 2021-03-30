package eod

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

const maxComboLength = 20

var combs = []string{
	"+",
	",",
}

func (b *EoD) cmdHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := b.newMsgNormal(m)
	rsp := b.newRespNormal(m)

	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	if strings.HasPrefix(m.Content, "*2") {
		if !b.checkServer(msg, rsp) {
			return
		}
		lock.RLock()
		dat, exists := b.dat[msg.GuildID]
		lock.RUnlock()
		if !exists {
			return
		}
		if dat.combCache == nil {
			dat.combCache = make(map[string]comb)
		}
		comb, exists := dat.combCache[msg.Author.ID]
		if !exists {
			return
		}
		if comb.elem3 != "" {
			b.combine([]string{comb.elem3}, msg, rsp)
			return
		}
		b.combine(comb.elems, msg, rsp)
		return
	}

	for _, comb := range combs {
		if strings.Contains(m.Content, comb) {
			if !b.checkServer(msg, rsp) {
				return
			}
			parts := strings.Split(m.Content, comb)
			if len(parts) < 2 {
				return
			}
			for i, part := range parts {
				parts[i] = strings.TrimSpace(part)
			}
			if len(parts) > maxComboLength {
				parts = parts[:maxComboLength]
			}
			set := make(map[string]empty, len(parts))
			for _, val := range parts {
				if len(val) > 240 {
					val = val[:240]
				}
				set[val] = empty{}
			}
			parts = make([]string, len(set))
			i := 0
			for k := range set {
				parts[i] = k
				i++
			}
			b.combine(parts, msg, rsp)
			return
		}
	}

	if strings.HasPrefix(m.Content, "?") {
		name := strings.TrimSpace(m.Content[1:])
		isGood := true
		if strings.Contains(name, "?") {
			for _, val := range name {
				if val != '?' {
					isGood = false
					break
				}
			}
		}
		if !isGood || len(name) == 0 {
			return
		}
		b.infoCmd(name, msg, rsp)
	}
}
