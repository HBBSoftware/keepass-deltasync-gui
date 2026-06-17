// SPDX-License-Identifier: GPL-3.0-or-later

package main

// i18n holder en meget enkel to-sprogs-ordbog (dansk/engelsk). GUI'en er
// bevidst en tynd skal oven på keepass-deltasync-CLI'en, så vi har kun brug
// for et lille sæt strenge. Default er dansk; brugeren kan skifte til engelsk
// i Indstillinger, og valget gemmes i gui.json sammen med stien til CLI'en.

type lang string

const (
	langDA lang = "da"
	langEN lang = "en"
)

// L er den aktive ordbog. Sættes ved opstart fra gemt indstilling og opdateres
// når brugeren skifter sprog.
var L = dicts[langDA]

func setLang(l lang) {
	if d, ok := dicts[l]; ok {
		L = d
	}
}

// dict er alle de strenge UI'en bruger. Felterne er navngivet efter hvor de
// optræder, så det er nemt at finde dem i koden.
type dict struct {
	AppTitle string

	// Generelt
	OK, Cancel, Close, Save, Back, Next, Finish, Browse string
	Working, Done, Error                                string
	Copy, DBInfoTitle                                   string

	// CLI-lokalisering
	CLINotFound     string
	CLILocate       string
	CLILocateHint   string
	CLIPathLabel    string
	CLISelectBinary string

	// Wizard
	WizardWelcomeTitle string
	WizardWelcomeBody  string
	WizardStart        string
	WizardStepEnroll   string
	WizardStepDB       string
	ServerURL          string
	ServerURLHint      string
	EnrollToken        string
	EnrollTokenHint    string
	DeviceName         string
	DeviceNameHint     string
	EnrollButton       string
	EnrollOK           string
	WizardAddDB        string
	WizardAddDBBody    string
	DBName             string
	DBNameHint         string
	KdbxFile           string
	KdbxFileHint       string
	CreateDBButton     string
	SkipForNow         string
	WizardDone         string
	WizardDoneBody     string

	// Dashboard
	TabDatabases string
	TabDevices   string
	TabActivity  string
	TabSettings  string
	ActivityHint string
	Clear        string

	// Enheder-fane
	ColEnrolled     string
	ColLastSeen     string
	ThisDevice      string
	DevCount        string // "%d"
	DeviceInfoTitle string

	AddDevice           string
	AddDeviceCreate     string
	RemoveDevice        string
	Username            string
	AdminToken          string
	EnrollTokenCreated  string
	SelectDeviceFirst   string
	CannotRemoveCurrent string
	ConfirmRemoveDevice string // "%s" = enhedsnavn
	StatusBox           string
	NotEnrolled         string
	Refresh             string
	Sync                string
	SyncSelected        string
	SyncAll             string
	SelectFirst         string
	AddDatabase         string
	ForgetDatabase      string
	ConfirmForget       string // "%s" = db-navn
	NoDatabases         string
	DBCount             string
	ColName             string
	ColStatus           string
	ColID               string
	ColCreated          string
	ColPath             string
	BoundLocally        string
	OnServerOnly        string

	// Detalje-panel: medlemmer/enheder koblet til den valgte database.
	MembersTitle       string
	MembersOf          string // "%s" = db-navn
	SelectToSeeMembers string
	MemberCount        string // "%d"
	MembersNeedBound   string
	MembersUnavailable string
	NoMembers          string
	ColRole            string
	ColUser            string
	ColDisplay         string
	ColAdded           string

	// Deling (share / unshare / init-shared)
	ShareDatabase     string
	ShareTitle        string // "%s" = db-navn
	ShareWith         string
	RemoveMember      string
	SelectMemberFirst string
	CannotRemoveOwner string
	ConfirmUnshare    string // "%s" = bruger, "%q" = db
	SetupShared       string
	SetupSharedTitle  string // "%s" = remote-navn
	AlreadyLocal      string
	NewLocalPassword  string
	BindExisting      string
	BoundNowSync      string
	MasterPwd         string
	MasterPwdFor      string
	Language          string
	ResetEnroll       string
}

var dicts = map[lang]*dict{
	langDA: {
		AppTitle: "KeePass Delta-Sync",

		OK: "OK", Cancel: "Annullér", Close: "Luk", Save: "Gem",
		Back: "Tilbage", Next: "Næste", Finish: "Færdig", Browse: "Vælg…",
		Working: "Arbejder…", Done: "Færdig", Error: "Fejl",
		Copy: "Kopiér", DBInfoTitle: "Database-info",

		CLINotFound:     "Kan ikke finde programmet 'keepass-deltasync'.",
		CLILocate:       "Find keepass-deltasync",
		CLILocateHint:   "Peg på keepass-deltasync-programmet (CLI'en). Det leveres normalt ved siden af denne app.",
		CLIPathLabel:    "Sti til CLI",
		CLISelectBinary: "Vælg keepass-deltasync-programmet",

		WizardWelcomeTitle: "Velkommen",
		WizardWelcomeBody: "Denne guide hjælper dig i gang med at synkronisere din KeePass-database " +
			"mellem dine enheder.\n\nDu skal bruge:\n  • Server-adressen fra din administrator\n  • Et enrollment-token (engangskode)\n  • Din .kdbx-fil",
		WizardStart:      "Start guide",
		WizardStepEnroll: "Trin 1 af 2 — Tilmeld enhed",
		WizardStepDB:     "Trin 2 af 2 — Tilføj database",
		ServerURL:        "Server-adresse",
		ServerURLHint:    "F.eks. https://deltasync.example.dk",
		EnrollToken:      "Enrollment-token",
		EnrollTokenHint:  "Engangskoden du fik af din administrator",
		DeviceName:       "Enhedsnavn (valgfrit)",
		DeviceNameHint:   "Vises i administrationen. Tomt = computerens navn.",
		EnrollButton:     "Tilmeld denne enhed",
		EnrollOK:         "Enheden er tilmeldt!",
		WizardAddDB:      "Tilføj din første database",
		WizardAddDBBody:  "Registrér en lokal .kdbx-fil til synkronisering. Du kan altid tilføje flere senere.",
		DBName:           "Navn",
		DBNameHint:       "Et kort navn, f.eks. 'privat' eller 'arbejde'",
		KdbxFile:         "KeePass-fil (.kdbx)",
		KdbxFileHint:     "Stien til din lokale database",
		CreateDBButton:   "Opret database",
		SkipForNow:       "Spring over",
		WizardDone:       "Alt klar!",
		WizardDoneBody:   "Du er klar til at synkronisere. Brug knappen 'Synkronisér' på en database for at sende og hente ændringer.",

		TabDatabases: "Databaser",
		TabDevices:   "Enheder",
		TabActivity:  "Aktivitet",
		TabSettings:  "Indstillinger",
		ActivityHint: "Output fra CLI-kald (sync, tilføj/fjern, fejl …) vises her.",
		Clear:        "Ryd",

		ColEnrolled:     "Tilmeldt",
		ColLastSeen:     "Sidst set",
		ThisDevice:      "● denne enhed",
		DevCount:        "%d enhed(er) på kontoen",
		DeviceInfoTitle: "Enheds-info",

		AddDevice:           "Tilføj enhed",
		AddDeviceCreate:     "Generér token",
		RemoveDevice:        "Fjern enhed",
		Username:            "Brugernavn",
		AdminToken:          "Admin-token",
		EnrollTokenCreated:  "Enrollment-token til ny enhed",
		SelectDeviceFirst:   "Vælg først en enhed i listen.",
		CannotRemoveCurrent: "Du kan ikke fjerne den enhed du bruger lige nu herfra.",
		ConfirmRemoveDevice: "Fjern (tilbagekald) enheden %q? Dens token bliver ugyldig.",
		StatusBox:           "Status",
		NotEnrolled:         "Ikke tilmeldt endnu.",
		Refresh:             "Opdatér",
		Sync:                "Synkronisér",
		SyncSelected:        "Synkronisér valgte",
		SyncAll:             "Synkronisér alle",
		SelectFirst:         "Vælg først en database i listen.",
		AddDatabase:         "Tilføj database",
		ForgetDatabase:      "Glem database",
		ConfirmForget:       "Glem den lokale binding for %q? Selve .kdbx-filen og databasen på serveren røres IKKE — kun koblingen i denne klient fjernes.",
		NoDatabases:         "Ingen databaser endnu. Klik 'Tilføj database' for at komme i gang.",
		DBCount:             "%d database(r) tilkoblet",
		ColName:             "Navn",
		ColStatus:           "Status",
		ColID:               "ID",
		ColCreated:          "Oprettet",
		ColPath:             "Lokal sti",
		BoundLocally:        "● klar",
		OnServerOnly:        "○ kun på server",

		MembersTitle:       "Koblet til databasen",
		MembersOf:          "Koblet til %q",
		SelectToSeeMembers: "Vælg en database ovenfor for at se hvem der er koblet til den.",
		MemberCount:        "%d medlem(mer)",
		MembersNeedBound:   "Databasen skal være sat op lokalt for at vise medlemmer.",
		MembersUnavailable: "Kan ikke hente medlemmer (kun ejeren kan se dette):",
		NoMembers:          "Kun dig — ikke delt med andre endnu.",
		ColRole:            "Rolle",
		ColUser:            "Brugernavn",
		ColDisplay:         "Visningsnavn",
		ColAdded:           "Tilføjet",

		ShareDatabase:     "Del database",
		ShareTitle:        "Del %q med en bruger",
		ShareWith:         "Del med (brugernavn)",
		RemoveMember:      "Fjern medlem",
		SelectMemberFirst: "Vælg først et medlem i listen.",
		CannotRemoveOwner: "Ejeren kan ikke fjernes.",
		ConfirmUnshare:    "Fjern %s fra %q?",
		SetupShared:       "Sæt op lokalt",
		SetupSharedTitle:  "Sæt den delte database %q op lokalt",
		AlreadyLocal:      "Databasen er allerede sat op lokalt.",
		NewLocalPassword:  "Nyt lokalt password",
		BindExisting:      "Forbind eksisterende .kdbx (egen database, ny enhed)",
		BoundNowSync:      "Forbundet! Klik på ⟳ (Synkronisér) for at hente entries.",
		MasterPwd:         "Masterpassword",
		MasterPwdFor:      "Masterpassword for",
		Language:          "Sprog",
		ResetEnroll:       "Nulstil tilmelding",
	},
	langEN: {
		AppTitle: "KeePass Delta-Sync",

		OK: "OK", Cancel: "Cancel", Close: "Close", Save: "Save",
		Back: "Back", Next: "Next", Finish: "Finish", Browse: "Browse…",
		Working: "Working…", Done: "Done", Error: "Error",
		Copy: "Copy", DBInfoTitle: "Database info",

		CLINotFound:     "Could not find the 'keepass-deltasync' program.",
		CLILocate:       "Locate keepass-deltasync",
		CLILocateHint:   "Point to the keepass-deltasync program (the CLI). It usually ships next to this app.",
		CLIPathLabel:    "Path to CLI",
		CLISelectBinary: "Select the keepass-deltasync program",

		WizardWelcomeTitle: "Welcome",
		WizardWelcomeBody: "This guide helps you start syncing your KeePass database " +
			"across your devices.\n\nYou will need:\n  • The server address from your administrator\n  • An enrollment token (one-time code)\n  • Your .kdbx file",
		WizardStart:      "Start guide",
		WizardStepEnroll: "Step 1 of 2 — Enroll device",
		WizardStepDB:     "Step 2 of 2 — Add database",
		ServerURL:        "Server address",
		ServerURLHint:    "e.g. https://deltasync.example.dk",
		EnrollToken:      "Enrollment token",
		EnrollTokenHint:  "The one-time code from your administrator",
		DeviceName:       "Device name (optional)",
		DeviceNameHint:   "Shown in admin listings. Empty = computer name.",
		EnrollButton:     "Enroll this device",
		EnrollOK:         "Device enrolled!",
		WizardAddDB:      "Add your first database",
		WizardAddDBBody:  "Register a local .kdbx file for syncing. You can always add more later.",
		DBName:           "Name",
		DBNameHint:       "A short name, e.g. 'personal' or 'work'",
		KdbxFile:         "KeePass file (.kdbx)",
		KdbxFileHint:     "Path to your local database",
		CreateDBButton:   "Create database",
		SkipForNow:       "Skip for now",
		WizardDone:       "All set!",
		WizardDoneBody:   "You're ready to sync. Use the 'Sync' button on a database to send and fetch changes.",

		TabDatabases: "Databases",
		TabDevices:   "Devices",
		TabActivity:  "Activity",
		TabSettings:  "Settings",
		ActivityHint: "Output from CLI calls (sync, add/remove, errors …) shows here.",
		Clear:        "Clear",

		ColEnrolled:     "Enrolled",
		ColLastSeen:     "Last seen",
		ThisDevice:      "● this device",
		DevCount:        "%d device(s) on the account",
		DeviceInfoTitle: "Device info",

		AddDevice:           "Add device",
		AddDeviceCreate:     "Generate token",
		RemoveDevice:        "Remove device",
		Username:            "Username",
		AdminToken:          "Admin token",
		EnrollTokenCreated:  "Enrollment token for new device",
		SelectDeviceFirst:   "Select a device in the list first.",
		CannotRemoveCurrent: "You can't remove the device you're currently using from here.",
		ConfirmRemoveDevice: "Remove (revoke) device %q? Its token will become invalid.",
		StatusBox:           "Status",
		NotEnrolled:         "Not enrolled yet.",
		Refresh:             "Refresh",
		Sync:                "Sync",
		SyncSelected:        "Sync selected",
		SyncAll:             "Sync all",
		SelectFirst:         "Select a database in the list first.",
		AddDatabase:         "Add database",
		ForgetDatabase:      "Forget database",
		ConfirmForget:       "Forget the local binding for %q? The .kdbx file and the database on the server are NOT touched — only the binding in this client is removed.",
		NoDatabases:         "No databases yet. Click 'Add database' to get started.",
		DBCount:             "%d database(s) connected",
		ColName:             "Name",
		ColStatus:           "Status",
		ColID:               "ID",
		ColCreated:          "Created",
		ColPath:             "Local path",
		BoundLocally:        "● ready",
		OnServerOnly:        "○ server only",

		MembersTitle:       "Connected to the database",
		MembersOf:          "Connected to %q",
		SelectToSeeMembers: "Select a database above to see who is connected to it.",
		MemberCount:        "%d member(s)",
		MembersNeedBound:   "The database must be set up locally to show members.",
		MembersUnavailable: "Cannot fetch members (only the owner can see this):",
		NoMembers:          "Just you — not shared with anyone yet.",
		ColRole:            "Role",
		ColUser:            "Username",
		ColDisplay:         "Display name",
		ColAdded:           "Added",

		ShareDatabase:     "Share database",
		ShareTitle:        "Share %q with a user",
		ShareWith:         "Share with (username)",
		RemoveMember:      "Remove member",
		SelectMemberFirst: "Select a member in the list first.",
		CannotRemoveOwner: "The owner can't be removed.",
		ConfirmUnshare:    "Remove %s from %q?",
		SetupShared:       "Set up locally",
		SetupSharedTitle:  "Set up the shared database %q locally",
		AlreadyLocal:      "The database is already set up locally.",
		NewLocalPassword:  "New local password",
		BindExisting:      "Connect existing .kdbx (your own database, new device)",
		BoundNowSync:      "Connected! Click ⟳ (Sync) to fetch entries.",
		MasterPwd:         "Master password",
		MasterPwdFor:      "Master password for",
		Language:          "Language",
		ResetEnroll:       "Reset enrollment",
	},
}
