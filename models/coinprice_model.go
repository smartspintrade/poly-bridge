/*
 * Copyright (C) 2020 The poly network Authors
 * This file is part of The poly network library.
 *
 * The  poly network  is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The  poly network  is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 * You should have received a copy of the GNU Lesser General Public License
 * along with The poly network .  If not, see <http://www.gnu.org/licenses/>.
 */

package models

import (
	"math/big"
)

const (
	TokenTypeErc20 uint8 = iota
	TokenTypeErc721
)

type TokenBasic struct {
	Name            string         `gorm:"primaryKey;size:64;not null"`
	Precision       uint64         `gorm:"type:bigint(20);not null"`
	Price           int64          `gorm:"size:64;not null"`
	Ind             uint64         `gorm:"type:bigint(20);not null"` // 显示价格是否可用
	Time            int64          `gorm:"type:bigint(20);not null"`
	Property        int64          `gorm:"type:bigint(20);not null"` // token是否上线, 1为上线
	Standard        uint8          `gorm:"type:int(8);not null"`     // 0为erc20， 1为erc721
	Meta            string         `gorm:"type:varchar(128)"`
	TotalAmount     *BigInt        `gorm:"type:varchar(64)"`
	TotalCount      uint64         `gorm:"type:bigint(20)"`
	StatsUpdateTime uint64         `gorm:"type:bigint(20)"`
	SocialTwitter   string         `gorm:"type:varchar(256)"`
	SocialTelegram  string         `gorm:"type:varchar(256)"`
	SocialWebsite   string         `gorm:"type:varchar(256)"`
	SocialOther     string         `gorm:"type:varchar(256)"`
	MetaFetcherType int            `gorm:"type:int(8);not null"` // nft meta profile fetcher type, e.g: unknown 0, opensea: 1, standard: 2,
	PriceMarkets    []*PriceMarket `gorm:"foreignKey:TokenBasicName;references:Name"`
	Tokens          []*Token       `gorm:"foreignKey:TokenBasicName;references:Name"`
}

type PriceMarket struct {
	TokenBasicName string      `gorm:"primaryKey;size:64;not null"`
	MarketName     string      `gorm:"primaryKey;size:64;not null"`
	Name           string      `gorm:"size:64;not null"`
	Price          int64       `gorm:"type:bigint(20);not null"`
	Ind            uint64      `gorm:"type:bigint(20);not null"`
	Time           int64       `gorm:"type:bigint(20);not null"`
	TokenBasic     *TokenBasic `gorm:"foreignKey:TokenBasicName;references:Name"`
}

type ChainFee struct {
	ChainId        uint64      `gorm:"primaryKey;type:bigint(20);not null"`
	TokenBasicName string      `gorm:"size:64;not null"`
	TokenBasic     *TokenBasic `gorm:"foreignKey:TokenBasicName;references:Name"`
	MaxFee         *BigInt     `gorm:"type:varchar(64);not null"`
	MinFee         *BigInt     `gorm:"type:varchar(64);not null"`
	ProxyFee       *BigInt     `gorm:"type:varchar(64);not null"`
	Ind            uint64      `gorm:"type:bigint(20);not null"`
	Time           int64       `gorm:"type:bigint(20);not null"`
}

type Token struct {
	Hash            string      `gorm:"primaryKey;size:66;not null"`
	ChainId         uint64      `gorm:"primaryKey;type:bigint(20);not null"`
	Name            string      `gorm:"size:64;not null"`
	Precision       uint64      `gorm:"type:bigint(20);not null"`
	TokenBasicName  string      `gorm:"size:64;not null"`
	Property        int64       `gorm:"type:bigint(20);not null"`
	Standard        uint8       `gorm:"type:int(8);not null"`
	AvailableAmount *BigInt     `gorm:"type:varchar(64)"`
	TokenBasic      *TokenBasic `gorm:"foreignKey:TokenBasicName;references:Name"`
	TokenMaps       []*TokenMap `gorm:"foreignKey:SrcTokenHash,SrcChainId;references:Hash,ChainId"`
}

type TokenMap struct {
	SrcChainId   uint64 `gorm:"primaryKey;type:bigint(20);not null"`
	SrcTokenHash string `gorm:"primaryKey;size:66;not null"`
	SrcToken     *Token `gorm:"foreignKey:SrcTokenHash,SrcChainId;references:Hash,ChainId"`
	DstChainId   uint64 `gorm:"primaryKey;type:bigint(20);not null"`
	DstTokenHash string `gorm:"primaryKey;size:66;not null"`
	DstToken     *Token `gorm:"foreignKey:DstTokenHash,DstChainId;references:Hash,ChainId"`
	Standard     uint8  `gorm:"type:int(8);not null"`
	Property     int64  `gorm:"type:bigint(20);not null"`
}

type WrapperTransactionWithToken struct {
	Hash         string  `gorm:"primaryKey;size:66;not null"`
	User         string  `gorm:"size:64"`
	SrcChainId   uint64  `gorm:"type:bigint(20);not null"`
	BlockHeight  uint64  `gorm:"type:bigint(20);not null"`
	Time         uint64  `gorm:"type:bigint(20);not null"`
	DstChainId   uint64  `gorm:"type:bigint(20);not null"`
	DstUser      string  `gorm:"type:varchar(66);not null"`
	ServerId     uint64  `gorm:"type:bigint(20);not null"`
	FeeTokenHash string  `gorm:"size:66;not null"`
	FeeToken     *Token  `gorm:"foreignKey:FeeTokenHash,SrcChainId;references:Hash,ChainId"`
	FeeAmount    *BigInt `gorm:"type:varchar(64);not null"`
	Status       uint64  `gorm:"type:bigint(20);not null"`
}

type CheckFee struct {
	ChainId     uint64
	Hash        string
	PayState    int
	Amount      *big.Float
	MinProxyFee *big.Float
}

type TimeStatistic struct {
	SrcChainId uint64 `gorm:"primaryKey;type:bigint(20);not null"`
	DstChainId uint64 `gorm:"primaryKey;type:bigint(20);not null"`
	Time       uint64 `gorm:"primaryKey;type:bigint(20);not null"`
}
