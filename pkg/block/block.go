/*
sysverity
Copyright (C) 2024  mutablefigment

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package block

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Block struct {
	ID                        int64
	IsGenisis                 bool
	Current_Block_Hash        string
	Previous_Block_Hash       string
	Signature_Bytes_Hex       string
	Genisis_Pubkey            string
	Boothash                  string
	UpdatedWithHashChainToken string
	UnixTimeStamp             string
}

func PrettyPrint(block *Block) {
	fmt.Printf("id %d \r\n", block.ID)
}

// TODO: add the sign hash function!

func HashGenisisBlock(block *Block) (string, error) {

	id := []byte{byte(block.ID)}
	boot_hash := []byte(block.Boothash)
	pubkey_bytes := []byte(block.Genisis_Pubkey)

	var all_bytes []byte
	all_bytes = append(all_bytes, id...)
	all_bytes = append(all_bytes, boot_hash...)
	all_bytes = append(all_bytes, pubkey_bytes...)

	gen_hash := sha256.Sum256(all_bytes)

	return hex.EncodeToString(gen_hash[:]), nil
}

func CacluatePrevBlockHash(block *Block) (string, error) {

	id := []byte{byte(block.ID)}
	prev_hash := []byte(block.Previous_Block_Hash)
	prev_sig := []byte(block.Signature_Bytes_Hex)

	var all_bytes []byte
	all_bytes = append(all_bytes, id...)
	all_bytes = append(all_bytes, prev_hash...)
	all_bytes = append(all_bytes, prev_sig...)

	prev_block_hash := sha256.Sum256(all_bytes)

	return hex.EncodeToString(prev_block_hash[:]), nil
}

func CacluateCurrBlockHash(block *Block) (string, error) {

	id := []byte{byte(block.ID)}
	prev_hash := []byte(block.Previous_Block_Hash)
	prev_sig := []byte(block.Signature_Bytes_Hex)

	var all_bytes []byte
	all_bytes = append(all_bytes, id...)
	all_bytes = append(all_bytes, prev_hash...)
	all_bytes = append(all_bytes, prev_sig...)

	prev_block_hash := sha256.Sum256(all_bytes)
	current_block_hash := hex.EncodeToString(prev_block_hash[:])

	fmt.Print(current_block_hash, "\r\n")
	fmt.Printf("id       : %x\r\n", id)
	fmt.Printf("prev_hash: %x\r\n", prev_hash)
	fmt.Printf("prev_sig : %x\r\n", prev_sig)
	fmt.Print(block)

	return current_block_hash, nil
}
