package internal

//go:generate gen-mysqlcrud -type Quiz

import (
	"time"
)

/*
	[ CRUD Tags Used By gen-mysqlcrud ]

	Pre-Requisites
		1) If RowTS type of column exists, must name it as RowTS time.Time
		2) If supporting logical delete, must have columns IsDeleted bool, DeletedUTC sql.NullTime
		3) Must have variable _originalValue interface{} `db:"-"`

	ggTable:"?"
	ggDel:"Logical" (Indicate Logical Delete Support)

	ggAutoCounter:"1,2,3,4,5,6,7,8,9"
		1 = Value: Auto Counter Name 															// (String)
		2 = - | Set-MID 																		// (null int64)
		3 = - | Set-TID 																		// (null int64)
		4 = - | Set-StoreID 																	// (null int64)
		5 = - | Set-ChainID 																	// (null int64)
		6 = - | Value: Prefix ATB Number If New 												// (null string, default 000000) if database has DBSetting, SettingEnum 9999 will contain DBNodeNumber if set during deploy, if so, ATB's BBBB is replaced with the DBNodeNumber during new counter creation
		7 = - | Value: Min Sequence Number If new 												// (null int64)
		8 = - | Value: Max Sequence Number If new 												// (null int64)
		9 = - | Value: Auto Reset Frequency If new 												// (null int32)

	ggAttr:"1,2,3,4,...9"
		1 = PK-Auto | PK-Custom | UK | FK | UK-FK | CUK1 | CUK2 | CUK1-FK | CUK2-FK | -			// column key type, - indicates N/A
		2 = FKTable.PKName (If #1 = FK) | De-UK (If #1 = UK) | -								// if #1 is FK, then #2 Must Be FK Table Name DOT PK Field Name, if #1 is UK, then #2 May Be 'De-UK' to avoid duplicate error, - indicates N/A
		3 = Null | -																			// column requirement, - indicates Required
		4 = # (Min Size or Range) | -															// column min size or range, - indicates no minimum
		5 = # (Max Size or Range) | -															// column max size or range, - indicates no maximum
		6 = Def-Ins | Def-Upd | Def-InsUpd | Set-Ins | Set-Upd | Set-InsUpd | Set-Del | -		// column value injected by database default in certain action, or set by caller in certain action, - indicates N/A
		7 = Check (On Data Changed) | -															// column value to be checked for changes, - indicates N/A
		8 = Init (Parameter) | Fill (Parameter) | -												// If column participates in Init() or Fill() functions, and act as its parameter, - indicates N/A
		9 = GetBy (Func) | GetRanged (Func) | GetPatterned (Func) | -							// create special function for this column, - indicates N/A (GetBy not for bool, float)
		10= F=x;y;z | -																			// used only with GetBy, GetRanged, GetPatterned defined in 9, F indicates filter parameter for Get action;
																								   note 1: if Get uses field itself with no additional filters, no need for F declare;
																								   note 2: x;y;z indicates additional fields as part of the filter for lookup (each separated by semi-colon)
																								   note 3: *x = * is to indicate x as optional parameter, meaning parameter is defined as pointer
																								   note 4: x?? = ?? can optional be placed at end of x, on data types integer, float, or date; so that the matching logic for x can be other than Equals
																										   x== (Equals) this is the default, no need to indicate as it is assumed
																										   x<> (Not Equals) overrides x to compare on Not Equal; x is field in db or in collection list, right side of <> is the input value to compare against
																										   x<< (Less Than) overrides x to compare on Less Than; x is field in db or in collection list, right side of << is the input value to compare against
																										   x>> (Greater Than) overrides x to compare on Greater Than; x is field in db or in collection list, right side of >> is the input value to compare against
																										   x<= (Less or Equal) overrides x to compare in Less Than or Equal To; x is field in db or in collection list, right side of <= is the input value to compare against
																									       x>= (Greater or Equal) overrides x to compare in Greater Than or Equal To; x is field in db or in collection list, right side of >= is the input value to compare against
																										   x>< (Between) overrides x to compare using BETWEEN; x is the field in db or in collection list, x is queried between xFrom and xTo inclusive as input values
		11= Enum=foldername.EnumTypeName | -												    // indicates if this field is an enum, and defines the foldername.EnumTypeName convention, such as menuchoice.MenuChoice (This tag may be omitted if not enum field)
*/

// Quiz record level dal model
type Quiz struct {
	ID int64 `db:"ID" ggAttr:"PK-Auto,-,-,-,-,-,-,-,GetBy" ggTable:"Quiz" ggDel:"None"`

	PacketNumber int       `db:"PacketNumber" ggAttr:"-,-,-,-,-,Set-InsUpd,Check,Fill,GetBy"`
	UnitNumber   int       `db:"UnitNumber" ggAttr:"-,-,-,-,-,Set-InsUpd,Check,Fill,-"`
	Difficulty   int       `db:"Difficulty" ggAttr:"-,-,-,-,-,Set-InsUpd,Check,Fill,-"`
	PDFLink      string    `db:"PDFLink" ggAttr:"-,-,-,-,-,Set-InsUpd,Check,Fill,-"`
	Timer        time.Time `db:"Timer" ggAttr:"-,-,-,-,-,Set-InsUpd,Check,Fill,-"`
	AdminID      int64     `db:"AdminID" ggAttr:"-,-,-,-,-,Set-InsUpd,Check,Fill,-"`
	Key1      string     `db:"Key1" ggAttr:"-,-,-,-,-,Set-InsUpd,Check,Fill,-"`
	Key2      string     `db:"Key2" ggAttr:"-,-,-,-,-,Set-InsUpd,Check,Fill,-"`
	Key3      string     `db:"Key3" ggAttr:"-,-,-,-,-,Set-InsUpd,Check,Fill,-"`

	// stores the original version of this struct, value set by Get or GetBy
	_originalValue interface{} `db:"-"`
}
