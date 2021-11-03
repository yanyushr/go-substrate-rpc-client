package types

import (
	"fmt"
	"hash"
	"strings"

	"github.com/yanyushr/go-substrate-rpc-client/v3/scale"
	"github.com/yanyushr/go-substrate-rpc-client/v3/xxhash"
)

type Param struct {
	Name    Text
	HasType Bool
	TypeId  UCompact
}

func (m *Param) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.Name)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.HasType)
	if err != nil {
		return err
	}

	if m.HasType {
		err = decoder.Decode(&m.TypeId)
		if err != nil {
			return err
		}
	}
	return nil
}

type Field struct {
	HasName     Bool
	Name        Text
	TypeId      UCompact
	HasTypeName Bool
	TypeName    Text
	Docs        []Text
}

func (f *Field) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&f.HasName)
	if err != nil {
		return err
	}

	if f.HasName {
		err = decoder.Decode(&f.Name)
		if err != nil {
			return err
		}
	}

	err = decoder.Decode(&f.TypeId)
	if err != nil {
		return err
	}

	err = decoder.Decode(&f.HasTypeName)
	if err != nil {
		return err
	}

	if f.HasTypeName {
		err = decoder.Decode(&f.TypeName)
		if err != nil {
			return err
		}
	}

	err = decoder.Decode(&f.Docs)
	if err != nil {
		return err
	}

	return nil
}

type TypeDefComposite struct {
	Fields []Field
}

type Variant struct {
	Name   Text
	Fields []Field
	Index  uint8
	Docs   []Text
}

type TypeDefVariant struct {
	Variants []Variant
}

type TypeDefSequence struct {
	TypeId UCompact
}

type TypeDefArray struct {
	Len    uint32
	TypeId UCompact
}

type TypeDefTuple struct {
	TypeIds []UCompact
}

type TypeDefPrimitive struct {
	Name Text
}

func (t *TypeDefPrimitive) Decode(decoder scale.Decoder) error {
	index, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch index {
	case 0:
		t.Name = "Bool"
	case 1:
		t.Name = "Char"
	case 2:
		t.Name = "String"
	case 3:
		t.Name = "U8"
	case 4:
		t.Name = "U16"
	case 5:
		t.Name = "U32"
	case 6:
		t.Name = "U64"
	case 7:
		t.Name = "U128"
	case 8:
		t.Name = "U256"
	case 9:
		t.Name = "I8"
	case 10:
		t.Name = "I16"
	case 11:
		t.Name = "I32"
	case 12:
		t.Name = "I64"
	case 13:
		t.Name = "I128"
	case 14:
		t.Name = "I256"
	}

	return nil
}

type TypeDefCompact struct {
	Compact UCompact
}

type TypeDefBitSequence struct {
	BitStoreType UCompact
	BitOrderType UCompact
}

type TypeDef struct {
	IsComposite          bool
	Composite            TypeDefComposite
	IsVariant            bool
	Variant              TypeDefVariant
	IsSequence           bool
	Sequence             TypeDefSequence
	IsArray              bool
	Array                TypeDefArray
	IsTuple              bool
	Tuple                TypeDefTuple
	IsPrimitive          bool
	Primitive            TypeDefPrimitive
	IsCompact            bool
	Compact              TypeDefCompact
	IsBitSequence        bool
	BitSequence          TypeDefBitSequence
	IsHistoricMetaCompat bool
	HistoricMetaCompat   string
}

func (t *TypeDef) Decode(decoder scale.Decoder) error {
	index, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch index {
	case 0:
		t.IsComposite = true
		err = decoder.Decode(&t.Composite)
		if err != nil {
			return err
		}
	case 1:
		t.IsVariant = true
		err = decoder.Decode(&t.Variant)
		if err != nil {
			return err
		}
	case 2:
		t.IsSequence = true
		err = decoder.Decode(&t.Sequence)
		if err != nil {
			return err
		}
	case 3:
		t.IsArray = true
		err = decoder.Decode(&t.Array)
		if err != nil {
			return err
		}
	case 4:
		t.IsTuple = true
		err = decoder.Decode(&t.Tuple)
		if err != nil {
			return err
		}
	case 5:
		t.IsPrimitive = true
		err = decoder.Decode(&t.Primitive)
		if err != nil {
			return err
		}
	case 6:
		t.IsCompact = true
		err = decoder.Decode(&t.Compact)
		if err != nil {
			return err
		}
	case 7:
		t.IsBitSequence = true
		err = decoder.Decode(&t.BitSequence)
		if err != nil {
			return err
		}
	case 8:
		t.IsHistoricMetaCompat = true
		err = decoder.Decode(&t.HistoricMetaCompat)
		if err != nil {
			return err
		}
	}
	return nil
}

type LookupType struct {
	Path   []Text
	Params []Param
	Def    TypeDef
	Docs   []Text
}

type PortableType struct {
	Id   UCompact
	Type LookupType
}

type PortableRegistry struct {
	Types []PortableType
}

type SignedExtensionMetadata struct {
	Identifier       Text
	Type             UCompact
	AdditionalSigned UCompact
}

type ExtrinsicV14 struct {
	Type             UCompact
	Version          U8
	SignedExtensions []SignedExtensionMetadata
}

type MetadataV14 struct {
	Lookup    PortableRegistry
	Pallets   []PalletMetadataV14
	Extrinsic ExtrinsicV14
	Types     UCompact
}

func (m *MetadataV14) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.Lookup)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Pallets)
	if err != nil {
		return err
	}
	err = decoder.Decode(&m.Extrinsic)
	if err != nil {
		return err
	}
	return decoder.Decode(&m.Types)
}

func (m MetadataV14) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(m.Pallets)
	if err != nil {
		return err
	}
	return encoder.Encode(m.Extrinsic)
}

func (m *MetadataV14) FindCallIndex(call string) (CallIndex, error) {
	s := strings.Split(call, ".")
	for _, pallet := range m.Pallets {
		if string(pallet.Name) != s[0] {
			continue
		}
		for _, tp := range m.Lookup.Types {
			if tp.Id.Cmp(pallet.Calls.TypeId) != 0 {
				continue
			} else {
				for _, method := range tp.Type.Def.Variant.Variants {
					if string(method.Name) != s[1] {
						continue
					}
					return CallIndex{pallet.Index, method.Index}, nil
				}
			}
		}

		return CallIndex{}, fmt.Errorf("method %v not found within module %v for call %v", s[1], pallet.Name, call)
	}

	return CallIndex{}, fmt.Errorf("module %v not found in metadata for call %v", s[0], call)
}

func (m *MetadataV14) FindEventNamesForEventID(eventID EventID) (Text, Text, error) {
	for _, mod := range m.Pallets {
		if !mod.HasEvents {
			continue
		}
		if mod.Index != eventID[0] {
			continue
		}

		typeId := mod.Events.TypeId.Int64()

		for _, tp := range m.Lookup.Types {
			if tp.Id.Int64() == typeId {
				for _, vars := range tp.Type.Def.Variant.Variants {
					if vars.Index == eventID[1] {
						return mod.Name, vars.Name, nil
					}
				}
			}
		}
	}
	return "", "", fmt.Errorf("module index %v out of range", eventID[0])
}

func (m *MetadataV14) FindStorageEntryMetadata(module string, fn string) (StorageEntryMetadata, error) {
	for _, mod := range m.Pallets {
		if !mod.HasStorage {
			continue
		}
		if string(mod.Storage.Prefix) != module {
			continue
		}
		for _, s := range mod.Storage.Items {
			if string(s.Name) != fn {
				continue
			}
			return s, nil
		}
		return nil, fmt.Errorf("storage %v not found within module %v", fn, module)
	}
	return nil, fmt.Errorf("module %v not found in metadata", module)
}

func (m *MetadataV14) FindConstantValue(module Text, constant Text) ([]byte, error) {
	for _, mod := range m.Pallets {
		if mod.Name == module {
			value, err := mod.FindConstantValue(constant)
			if err == nil {
				return value, nil
			}
		}
	}
	return nil, fmt.Errorf("could not find constant %s.%s", module, constant)
}

func (m *MetadataV14) ExistsModuleMetadata(module string) bool {
	for _, mod := range m.Pallets {
		if string(mod.Name) == module {
			return true
		}
	}
	return false
}

// type PortableRegistry struct {
// 	//todo implement SiType/Si1Type https://github.com/polkadot-js/api/blob/2801088b0a05e6bc505c2c449f1eddb31b15587d/packages/types/src/interfaces/scaleInfo/v1.ts#L25
// 	Types []byte
// }

type FunctionCallMetadataV14 struct {
	TypeId UCompact
}

type EventMetadataV14 struct {
	TypeId UCompact
}

type ErrorMetadataV14 struct {
	TypeId UCompact
}

type ModuleConstantMetadataV14 struct {
	Name          Text
	TypeId        UCompact
	Value         Bytes
	Documentation []Text
}

type PalletMetadataV14 struct {
	Name       Text
	HasStorage bool
	Storage    StorageMetadataV14
	HasCalls   bool
	Calls      FunctionCallMetadataV14
	HasEvents  bool
	Events     EventMetadataV14
	Constants  []ModuleConstantMetadataV14
	HasErrors  bool
	Errors     ErrorMetadataV14
	Index      uint8
}

func (m *PalletMetadataV14) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.Name)
	if err != nil {
		return err
	}
	err = decoder.Decode(&m.HasStorage)
	if err != nil {
		return err
	}

	if m.HasStorage {
		err = decoder.Decode(&m.Storage)
		if err != nil {
			return err
		}
	}

	err = decoder.Decode(&m.HasCalls)
	if err != nil {
		return err
	}

	if m.HasCalls {
		err = decoder.Decode(&m.Calls)
		if err != nil {
			return err
		}
	}

	err = decoder.Decode(&m.HasEvents)
	if err != nil {
		return err
	}

	if m.HasEvents {
		err = decoder.Decode(&m.Events)
		if err != nil {
			return err
		}
	}

	err = decoder.Decode(&m.Constants)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.HasErrors)
	if err != nil {
		return err
	}

	if m.HasErrors {
		err = decoder.Decode(&m.Errors)
		if err != nil {
			return err
		}
	}

	err = decoder.Decode(&m.Index)
	if err != nil {
		return err
	}
	return err
}

func (m PalletMetadataV14) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(m.Name)
	if err != nil {
		return err
	}

	err = encoder.Encode(m.HasStorage)
	if err != nil {
		return err
	}

	if m.HasStorage {
		err = encoder.Encode(m.Storage)
		if err != nil {
			return err
		}
	}

	err = encoder.Encode(m.HasCalls)
	if err != nil {
		return err
	}

	if m.HasCalls {
		err = encoder.Encode(m.Calls)
		if err != nil {
			return err
		}
	}

	err = encoder.Encode(m.HasEvents)
	if err != nil {
		return err
	}

	if m.HasEvents {
		err = encoder.Encode(m.Events)
		if err != nil {
			return err
		}
	}

	err = encoder.Encode(m.Constants)
	if err != nil {
		return err
	}

	err = encoder.Encode(m.Errors)
	if err != nil {
		return err
	}

	return encoder.Encode(m.Index)
}

func (m *PalletMetadataV14) FindConstantValue(constant Text) ([]byte, error) {
	for _, cons := range m.Constants {
		if cons.Name == constant {
			return cons.Value, nil
		}
	}
	return nil, fmt.Errorf("could not find constant %s", constant)
}

type StorageMetadataV14 struct {
	Prefix Text
	Items  []StorageFunctionMetadataV14
}

type MapTypeV14 struct {
	Hashers []StorageHasherV10
	Keys    UCompact
	Value   UCompact
}

type StorageFunctionMetadataV14 struct {
	Name          Text
	Modifier      StorageFunctionModifierV0
	Type          StorageFunctionTypeV14
	Fallback      Bytes
	Documentation []Text
}

func (s StorageFunctionMetadataV14) IsPlain() bool {
	return s.Type.IsType
}

func (s StorageFunctionMetadataV14) IsMap() bool {
	return s.Type.IsMap
}

func (s StorageFunctionMetadataV14) IsDoubleMap() bool {
	return s.Type.IsDoubleMap
}

func (s StorageFunctionMetadataV14) IsNMap() bool {
	return s.Type.IsNMap
}

func (s StorageFunctionMetadataV14) Hasher() (hash.Hash, error) {
	if s.Type.IsMap {
		return s.Type.AsMap.Hashers[0].HashFunc()
	}
	if s.Type.IsDoubleMap {
		return s.Type.AsDoubleMap.Hasher.HashFunc()
	}
	if s.Type.IsNMap {
		return nil, fmt.Errorf("only Map and DoubleMap have a Hasher")
	}
	return xxhash.New128(nil), nil
}

func (s StorageFunctionMetadataV14) Hasher2() (hash.Hash, error) {
	if !s.Type.IsDoubleMap {
		return nil, fmt.Errorf("only DoubleMaps have a Hasher2")
	}
	return s.Type.AsDoubleMap.Key2Hasher.HashFunc()
}

func (s StorageFunctionMetadataV14) Hashers() ([]hash.Hash, error) {
	if !s.Type.IsNMap {
		return nil, fmt.Errorf("only NMaps have Hashers")
	}

	hashers := make([]hash.Hash, len(s.Type.AsNMap.Hashers))
	for i, hasher := range s.Type.AsNMap.Hashers {
		hasherFn, err := hasher.HashFunc()
		if err != nil {
			return nil, err
		}
		hashers[i] = hasherFn
	}
	return hashers, nil
}

type StorageFunctionTypeV14 struct {
	IsType      bool
	AsType      UCompact // 0
	IsMap       bool
	AsMap       MapTypeV14 // 1
	IsDoubleMap bool
	AsDoubleMap DoubleMapTypeV10 // 2
	IsNMap      bool
	AsNMap      NMapTypeV13 // 3
}

func (s *StorageFunctionTypeV14) Decode(decoder scale.Decoder) error {
	var t uint8
	err := decoder.Decode(&t)
	if err != nil {
		return err
	}

	switch t {
	case 0:
		s.IsType = true
		err = decoder.Decode(&s.AsType)
		if err != nil {
			return err
		}
	case 1:
		s.IsMap = true
		err = decoder.Decode(&s.AsMap)
		if err != nil {
			return err
		}
	case 2:
		s.IsDoubleMap = true
		err = decoder.Decode(&s.AsDoubleMap)
		if err != nil {
			return err
		}
	case 3:
		s.IsNMap = true
		err = decoder.Decode(&s.AsNMap)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("received unexpected type %v", t)
	}
	return nil
}

func (s StorageFunctionTypeV14) Encode(encoder scale.Encoder) error {
	switch {
	case s.IsType:
		err := encoder.PushByte(0)
		if err != nil {
			return err
		}
		err = encoder.Encode(s.AsType)
		if err != nil {
			return err
		}
	case s.IsMap:
		err := encoder.PushByte(1)
		if err != nil {
			return err
		}
		err = encoder.Encode(s.AsMap)
		if err != nil {
			return err
		}
	case s.IsDoubleMap:
		err := encoder.PushByte(2)
		if err != nil {
			return err
		}
		err = encoder.Encode(s.AsDoubleMap)
		if err != nil {
			return err
		}
	case s.IsNMap:
		err := encoder.PushByte(3)
		if err != nil {
			return err
		}
		err = encoder.Encode(s.AsNMap)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("expected to be either type, map, double map or nmap but none was set: %v", s)
	}
	return nil
}

type NMapTypeV14 struct {
	Keys    []Type
	Hashers []StorageHasherV10
	Value   Type
}
