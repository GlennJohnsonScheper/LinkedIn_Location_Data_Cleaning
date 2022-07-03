# LinkedIn_Location_Data_Cleaning

You know how annoying it is to search your candidate database
when it contains very diverse human entered location strings.

Consider just a few pertaining to New York City:

		"New York, NY", "100 church street new york",
		"New York, NY", "120 Park Ave, New York, NY",
		"New York, NY", "150 East 42nd Street, New York, NY 10017",
		"New York, NY", "787 7th Ave, New York, NY",
		"New York, NY", "826 Broadway, New York, NY",
		"New York, NY", "Elmira, New York Area",
		"New York, NY", "Greater New York Area",
		"New York, NY", "Greater New York City Area ",
		"New York, NY", "Greater New York City Area , San Francisco Bay Area",
		"New York, NY", "Greater New York City Area and remote",
		"New York, NY", "Greater New York City Area",
		"New York, NY", "Greater New York City Area, USA",
		"New York, NY", "Greater New York City Area/Austin TX",
		"New York, NY", "Großraum New York City und Umgebung",
		"New York, NY", "NYC",
		"New York, NY", "NYC, NY",
		"New York, NY", "NYC, NewYork",
		"New York, NY", "NYC, Ny",
		"New York, NY", "New York City and San Francisco",
		"New York, NY", "New York City",
		"New York, NY", "New York City, NY",
		"New York, NY", "New York City, New York",
		"New York, NY", "New York City, New York, USA",
		"New York, NY", "New York City, Stati Uniti",
		"New York, NY", "New York",
		"New York, NY", "New York, NY / Nationwide",
		"New York, NY", "New York, NY",
		"New York, NY", "New York, New York",
		"New York, NY", "New York,NY",
		"New York, NY", "New york",
		"New York, NY", "New york, NY",
		"New York, NY", "New york,ny",
		"New York, NY", "NewYork",
		"New York, NY", "Utica, New York Area",
		"New York, NY", "greater new york city area",

I still see two flies in the ointment just in those few lines.
Anyway, pareto principle, good enough for unpaid own business.

This .GO code shows some heuristics how I approached the task,
and the generated C# code map table may be valuable in itself.
