// C:\A\repo\LinkedIn_Location_Data_Cleaning\LinkedIn_Location_Data_Cleaning.go

// Clean-up user-entered LinkedIn locations.
// Lossy, cluster what can, discard bizarre.

// After all this coding, I still found one ringer in output to modify:
// 		"Shelton, CT", "1 Basking Ridge", // one manual fix here

// Note. I forgot to let the idempotent perfect location string map to itself.
// So the eventual C# program that will use this generated code table should
// add such to its dictionary of string to string mapping raw to clean text.

// Doing that, the next C# program fixed/passed 19K person top location names
// with ~300 rejects, which were almost all due to naming foreign countries.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// String-2-Sort by descending len():
type S2S []string

func (s2s S2S) Len() int           { return len(s2s) }
func (s2s S2S) Swap(a, b int)      { s2s[a], s2s[b] = s2s[b], s2s[a] }
func (s2s S2S) Less(a, b int) bool { return len(s2s[a]) > len(s2s[b]) }

// Ignoring "Alphaville", a fictional city in The Sims Online game.
// Omitting cities with an apostraphe in name. But recovered later!

// var perfected = []string{...
var perfected = S2S{
	"Aberdeen, MD",
	"Acton, MA",
	"Acworth, GA",
	"Addison, TX",
	"Adelphi, MD",
	"Agoura Hills, CA",
	"Akron, OH",
	"Alameda, CA",
	"Alamo, CA",
	"Albany, NY",
	"Albuquerque, NM",
	"Alexander, AR",
	"Alexandria, VA",
	"Algonquin, IL",
	"Alhambra, CA",
	"Aliso Viejo, CA",
	"Allen, TX",
	"Allentown, PA",
	"Alpharetta, GA",
	"Altamonte Springs, FL",
	"Alviso, CA",
	"American Canyon, CA",
	"American Fork, UT",
	"Ames, IA",
	"Amherst, MA",
	"Anaheim, CA",
	"Anchorage, AK",
	"Andover, MA",
	"Ann Arbor, MI",
	"Annandale, VA",
	"Annapolis Junction, MD",
	"Annapolis, MD",
	"Apex, NC",
	"Apopka, FL",
	"Apple Valley, CA",
	"Arcadia, CA",
	"Argyle, TX",
	"Arlington Heights, IL",
	"Arlington, VA",
	"Armonk, NY",
	"Arnold, MD",
	"Arvada, CO",
	"Ashburn, VA",
	"Ashland, MA",
	"Astoria, NY",
	"Athens, OH",
	"Atlanta, GA",
	"Auburn Hills, MI",
	"Auburn, WA",
	"Auburndale, MA",
	"Augusta, ME",
	"Aurora, IL",
	"Austin, TX",
	"Avenel, NJ",
	"Avon, IN",
	"Ayer, MA",
	"Babylon, NY",
	"Bainbridge Island, WA",
	"Baldwin Park, CA",
	"Baldwin, NY",
	"Ballston Lake, NY",
	"Baltimore, OH",
	"Barrington, NH",
	"Bartlett, IL",
	"Basking Ridge, NJ",
	"Baton Rouge, LA",
	"Battery Park, NY",
	"Bay Area, CA",
	"Bay Center, WA",
	"Bay Shore, NY",
	"Bayside, NY",
	"Beacon Falls, CT",
	"Beaumont, TX",
	"Beaverton, OR",
	"Bedford Hills, NY",
	"Bedford, MA",
	"Bedminster, NJ",
	"Bel Air, MD",
	"Bellaire, TX",
	"Belle Mead, NJ",
	"Bellevue, WA",
	"Bellingham, MA",
	"Belmont, CA",
	"Belvedere Tiburon, CA",
	"Bend, OR",
	"Benicia, CA",
	"Bennett, CO",
	"Bensalem, PA",
	"Bensenville, IL",
	"Benson, NC",
	"Bentonville, AR",
	"Berkeley Heights, NJ",
	"Berkeley, CA",
	"Berlin, MA",
	"Bethesda, MD",
	"Bethlehem, PA",
	"Beverly Hills, CA",
	"Biloxi, MS",
	"Binghamton, NY",
	"Birmingham, AL",
	"Bisbee, AZ",
	"Blacksburg, VA",
	"Bloomfield, CT",
	"Bloomfield, NJ",
	"Bloomington, IN",
	"Bloomington, MN",
	"Bloomsburg, PA",
	"Blue Bell, PA",
	"Boaz, AL",
	"Boca Raton, FL",
	"Bohemia, NY",
	"Boise, ID",
	"Bolingbrook, IL",
	"Bonita, CA",
	"Bonsall, CA",
	"Boonsboro, MD",
	"Boonton, NJ",
	"Boston, MA",
	"Bothell, WA",
	"Bouder, CO",
	"Boulder, CO",
	"Bound Brook, NJ",
	"Bowie, MD",
	"Boxborough, MA",
	"Boynton Beach, FL",
	"Bozeman, MT",
	"Brainerd, MN",
	"Braintree, MA",
	"Branford, CT",
	"Brea, CA",
	"Breinigsville, PA",
	"Brentwood, CA",
	"Briarcliff Manor, NY",
	"Bridgeport, CT",
	"Bridgewater, NJ",
	"Brighton, MA",
	"Brisbane, CA",
	"Bristol, CT",
	"Bristow, VA",
	"Broken Arrow, OK",
	"Bronx, NY",
	"Brookfield, WI",
	"Brookhaven, NY",
	"Brookline, MA",
	"Brooklyn, NY",
	"Broomfield, CO",
	"Buckeye, AZ",
	"Buffalo Grove, IL",
	"Buffalo, NY",
	"Buford, GA",
	"Burbank, CA",
	"Burleson, TX",
	"Burlingame, CA",
	"Burlington, MA",
	"Burnsville, MN",
	"Calabasas, CA",
	"Camarillo, CA",
	"Cambridge, MA",
	"Camden, NJ",
	"Camp Hill, PA",
	"Campbell, CA",
	"Canandaigua, NY",
	"Canoga Park, CA",
	"Canonsburg, PA",
	"Canton, OH",
	"Cape Girardeau, MO",
	"Capitol Heights, MD",
	"Cardiff By The Sea, CA",
	"Carefree, AZ",
	"Carle Place, NY",
	"Carlsbad, CA",
	"Carmel, IN",
	"Carpinteria, CA",
	"Carrollton, TX",
	"Carson City, NV",
	"Cary, NC",
	"Castle Rock, CO",
	"Castro Valley, CA",
	"Catharpin, VA",
	"Cave Creek, AZ",
	"Cedar Creek, TX",
	"Cedar Park, TX",
	"Cedar Rapids, IA",
	"Centennial, CO",
	"Center Valley, PA",
	"Centreville, VA",
	"Century City, CA",
	"Cerritos, CA",
	"Chambersburg, PA",
	"Champaign, IL",
	"Chandler, AZ",
	"Chanhassen, MN",
	"Chantilly, VA",
	"Chapel Hill, NC",
	"Charleston, SC",
	"Charlestown, MA",
	"Charlotte, NC",
	"Charlottesville, VA",
	"Chatsworth, CA",
	"Chattanooga, TN",
	"Chelmsford, MA",
	"Cherry Hill, NJ",
	"Chesapeake Beach, MD",
	"Chesterbrook, PA",
	"Chestnut Hill, MA",
	"Chevy Chase, MD",
	"Cheyenne, WY",
	"Chicago, IL",
	"Chico, CA",
	"Chino Hills, CA",
	"Chula Vista, CA",
	"Cincinnati, OH",
	"Clara, CA",
	"Claremont, CA",
	"Clark, NJ",
	"Clarksburg, MD",
	"Clarksville, VA",
	"Clayton, MO",
	"Clearwater Tampa, FL",
	"Clearwater, FL",
	"Cleveland, OH",
	"Cliffside Park, NJ",
	"Clifton, VA",
	"Clinton, MA",
	"Coalport, PA",
	"Cocoa, FL",
	"College Park, MD",
	"College Station, TX",
	"Collegeville, PA",
	"Colleyville, TX",
	"Collierville, TN",
	"Colonia, NJ",
	"Colorado Springs, CO",
	"Columbia, MD",
	"Columbus, OH",
	"Concord, NC",
	"Conroe, TX",
	"Conshohocken, PA",
	"Coppell, TX",
	"Coralville, IA",
	"Coraopolis, PA",
	"Corona, CA",
	"Corte Madera, CA",
	"Cortlandt Manor, NY",
	"Corvallis, OR",
	"Costa Mesa, CA",
	"Council Bluffs, IA",
	"Covina, CA",
	"Covington, GA",
	"Covington, KY",
	"Cranbury, NJ",
	"Cranford, NJ",
	"Cresskill, NJ",
	"Creston, CA",
	"Croton On Hudson, NY",
	"Crown Point, IN",
	"Culver City, CA",
	"Cumming, GA",
	"Cupertino, CA",
	"Dallas, TX",
	"Daly City, CA",
	"Damascus, MD",
	"Danbury, CT",
	"Danvers, MA",
	"Danville, CA",
	"Danville, KY",
	"Davenport, IA",
	"Davidsonville, MD",
	"Davis, CA",
	"Dayton, NJ",
	"Dayton, OH",
	"Dearborn, MI",
	"Dedham, MA",
	"Deerfield, IL",
	"Dekalb, IL",
	"Del Mar, CA",
	"Delray Beach, FL",
	"Deltona, FL",
	"Denton, TX",
	"Denver, CO",
	"Denville, NJ",
	"Des Moines, IA",
	"Des Plaines, IL",
	"Desert Hot Springs, CA",
	"Desert Oasis, AZ",
	"Detroit, MI",
	"Devon, PA",
	"Diamond Bar, CA",
	"Dillsburg, PA",
	"Dobbs Ferry, NY",
	"Dover, NH",
	"Downers Grove, IL",
	"Downey, CA",
	"Downingtown, PA",
	"Doylestown, PA",
	"Draper, UT",
	"Dresher, PA",
	"Duarte, CA",
	"Dublin, CA",
	"Dulles, VA",
	"Duluth, GA",
	"Dunwoody, GA",
	"Durham, NC",
	"Duvall, WA",
	"Eagan, MN",
	"East Brunswick, NJ",
	"East Elmhurst, NY",
	"East Freetown, MA",
	"East Greenwich, RI",
	"East Lansing, MI",
	"East Rutherford, NJ",
	"East Stroudsburg, PA",
	"Eastgate, WA",
	"Easton, CT",
	"Eatontown, NJ",
	"Eden Praire, MN",
	"Eden Prairie, MN",
	"Edgewater, NJ",
	"Edinburg, TX",
	"Edison, NJ",
	"Edmonds, WA",
	"El Cajon, CA",
	"El Cerrito, CA",
	"El Dorado Hills, CA",
	"El Monte, CA",
	"El Paso, TX",
	"El Segundo, CA",
	"Elgin, IL",
	"Elizabeth, NJ",
	"Elizabethtown, PA",
	"Elk Grove, CA",
	"Elkridge, MD",
	"Ellensburg, WA",
	"Ellicott City, MD",
	"Elmont, NY",
	"Elmsford, NY",
	"Elmwood Park, NJ",
	"Emeryville, CA",
	"Encinitas, CA",
	"Encino, CA",
	"Englewood Cliffs, NJ",
	"Englewood, CO",
	"Erie, CO",
	"Escondido, CA",
	"Estero, FL",
	"Eugene, OR",
	"Eureka, MO",
	"Evanston, IL",
	"Evansville, IN",
	"Everett, WA",
	"Ewing, NJ",
	"Exton, PA",
	"Fair Lawn, NJ",
	"Fair Oaks, VA",
	"Fairfax Station, VA",
	"Fairfax, VA",
	"Fairfield, IA",
	"Fairmont, WV",
	"Fairview, NJ",
	"Fall City, WA",
	"Fallbrook, CA",
	"Falls Church, VA",
	"Fanwood, NJ",
	"Fargo, ND",
	"Farmingdale, NJ",
	"Farmington, MI",
	"Fayetteville, NC",
	"Federal Way, WA",
	"Finksburg, MD",
	"Firestone, CO",
	"Fishers, IN",
	"Flagstaff, AZ",
	"Flemington, NJ",
	"Floresville, TX",
	"Florham Park, NJ",
	"Flower Mound, TX",
	"Flushing, NY",
	"Folsom, CA",
	"Fontana, CA",
	"Foothill Ranch, CA",
	"Fords, NJ",
	"Forest Hill, MD",
	"Forest Hills, NY",
	"Foristell, MO",
	"Fort Bragg, NC",
	"Fort Collins, CO",
	"Fort Huachuca, AZ",
	"Fort Lauderdale, FL",
	"Fort Mill, SC",
	"Fort Myers, FL",
	"Fort Smith, AR",
	"Fort Walton Beach, FL",
	"Fort Wayne, IN",
	"Fort Worth, TX",
	"Fortville, IN",
	"Foster City, CA",
	"Fountain Hills, AZ",
	"Fountain Valley, CA",
	"Framingham, MA",
	"Franklin Lakes, NJ",
	"Franklin Park, NJ",
	"Franklin Square, NY",
	"Franktown, CO",
	"Frederick, CO",
	"Fredericksburg, VA",
	"Freehold, NJ",
	"Freeport, NY",
	"Fremont, CA",
	"Fresno, CA",
	"Frisco, TX",
	"Fullerton, CA",
	"Fulton, MD",
	"Fuquay Varina, NC",
	"Gainesville, VA",
	"Gaithersburg, MD",
	"Gambrills, MD",
	"Garden Grove, CA",
	"Gardiner, NY",
	"Geneva, IL",
	"Georgetown, TX",
	"Germantown, MD",
	"Gerrardstown, WV",
	"Gibsonville, NC",
	"Gig Harbor, WA",
	"Gilbert, AZ",
	"Gilcrest, CO",
	"Gilroy, CA",
	"Glen Allen, VA",
	"Glen Gardner, NJ",
	"Glen Oaks, NY",
	"Glendale, CA",
	"Glendora, CA",
	"Glenside, PA",
	"Glenview, IL",
	"Glenwood, MD",
	"Golden, CO",
	"Goleta, CA",
	"Goodyear, AZ",
	"Goose Creek, SC",
	"Granada Hills, CA",
	"Grand Forks, ND",
	"Grand Rapids, MI",
	"Granite Bay, CA",
	"Great Barrington, MA",
	"Great Falls, VA",
	"Great Neck, NY",
	"Greeley, CO",
	"Green Bay, WI",
	"Greencastle, PA",
	"Greeneville, TN",
	"Greenville, WI",
	"Greenwich, CT",
	"Greenwood Village, CO",
	"Hacienda Heights, CA",
	"Hackensack, NJ",
	"Half Moon Bay, CA",
	"Hamden, CT",
	"Hammond, IN",
	"Harpers Ferry, WV",
	"Harrisburg, PA",
	"Harrison, NJ",
	"Harrison, NY",
	"Harrisonburg, VA",
	"Hartford, CT",
	"Harvard, MA",
	"Hauppauge, NY",
	"Hawthorne, NY",
	"Hayden, ID",
	"Haydenville, MA",
	"Haymarket, VA",
	"Hayward, CA",
	"Hazlet, NJ",
	"Heathrow, FL",
	"Hempstead, NY",
	"Henderson, TN",
	"Hercules, CA",
	"Herndon, VA",
	"Herriman, UT",
	"Hesperia, CA",
	"Hicksville, NY",
	"High Point, NC",
	"Highland Park, NJ",
	"Highlands Ranch, CO",
	"Highpoint, NC",
	"Hightstown, NJ",
	"Hilliard, OH",
	"Hillsboro, OR",
	"Hillsborough, NJ",
	"Hillsdale, NJ",
	"Hilton, NY",
	"Hingham, MA",
	"Ho Ho Kus, NJ",
	"Hoboken, NJ",
	"Hoffman Estates, IL",
	"Holden, MA",
	"Holliston, MA",
	"Holly Springs, NC",
	"Hollywood, FL",
	"Holmdel, NJ",
	"Holt, MI",
	"Holtsville, NY",
	"Holyoke, MA",
	"Honolulu, HI",
	"Hood River, OR",
	"Hopewell Junction, NY",
	"Hopewell, NJ",
	"Hopewell, VA",
	"Hopkinton, MA",
	"Horseshoe Bay, TX",
	"Horsham, PA",
	"Houghton, MI",
	"Houston, TX",
	"Huachuca City, AZ",
	"Hudson, NH",
	"Huntersville, NC",
	"Huntington Beach, CA",
	"Huntington Station, NY",
	"Huntington, WV",
	"Huntsville, AL",
	"Hyattsville, MD",
	"Idaho Falls, ID",
	"Imperial Beach, CA",
	"Independence, MO",
	"Indianapolis, IN",
	"Indianola, IA",
	"Indian River, MI",
	"Indian Trail, NC",
	"Iowa City, IA",
	"Irvine, CA",
	"Irving, TX",
	"Irvington, NY",
	"Irwin, PA",
	"Irwindale, CA",
	"Iselin, NJ",
	"Islandia, NY",
	"Issaquah, WA",
	"Jackson Heights, NY",
	"Jackson, NJ",
	"Jackson, WY",
	"Jacksonville, FL",
	"Jamestown, OH",
	"Jasper, AL",
	"Jericho, NY",
	"Jersey City, NJ",
	"Johns Creek, GA",
	"Johnson City, TN",
	"Jolla, CA",
	"Jonesboro, AR",
	"Julian, CA",
	"Juno Beach, FL",
	"Jupiter, FL",
	"Kalamazoo, MI",
	"Kansas City, MO",
	"Kearney, NE",
	"Kearny, NJ",
	"Keller, TX",
	"Kendall Park, NJ",
	"Kenmore, WA",
	"Kennebunk, ME",
	"Kensington, MD",
	"Kent, WA",
	"Kentwood, MI",
	"King Of Prussia, PA",
	"Kingwood, TX",
	"Kirkland, WA",
	"Kissimmee, FL",
	"Knoxville, TN",
	"La Crescenta, CA",
	"La Crosse, WI",
	"La Jolla, CA",
	"La Mesa, CA",
	"La Mirada, CA",
	"La Palma, CA",
	"La Verne, CA",
	"Ladera Ranch, CA",
	"Lafayette, LA",
	"Laguna Niguel, CA",
	"Lake Elmo, MN",
	"Lake Elsinore, CA",
	"Lake Forest, IL",
	"Lake Grove, NY",
	"Lake Mary, FL",
	"Lake Oswego, OR",
	"Lake Stevens, WA",
	"Lake Success, NY",
	"Lake Worth, FL",
	"Lakeland, FL",
	"Lakeside, CA",
	"Lakeville, MN",
	"Lakewood, CA",
	"Lancaster, CA",
	"Lanham, MD",
	"Lansdale, PA",
	"Lansing, MI",
	"Larkspur, CA",
	"Las Vegas, NV",
	"Laurel, MS",
	"Laveen, AZ",
	"Lavergne, TN",
	"Lawrence, KS",
	"Lawrenceburg, IN",
	"Lawrenceville, GA",
	"Leander, TX",
	"Lebanon, NJ",
	"Lees Summit, MO",
	"Leesburg, VA",
	"Lehi, UT",
	"Lenexa, KS",
	"Lewis Center, OH",
	"Lexington, SC",
	"Liberty Lake, WA",
	"Liberty Township, OH",
	"Lilburn, GA",
	"Lima, OH",
	"Lincoln, NE",
	"Lincolnshire, IL",
	"Lindon, UT",
	"Litchfield Park, AZ",
	"Lithia Springs, GA",
	"Little Chute, WI",
	"Little Elm, TX",
	"Little Rock, AR",
	"Little Silver, NJ",
	"Littleton, MA",
	"Livermore, CA",
	"Lodi, NJ",
	"Logan, UT",
	"Lompoc, CA",
	"Long Beach, CA",
	"Long Island City, NY",
	"Long Island, NY",
	"Longmont, CO",
	"Loomis, CA",
	"Lorton, VA",
	"Los Alamos, NM",
	"Los Altos, CA",
	"Los Angeles, CA",
	"Los Gatos, CA",
	"Louisville, CO",
	"Louisville, KY",
	"Loveland, CO",
	"Lowell, MA",
	"Lubbock, TX",
	"Lumberton, NJ",
	"Lutherville Timonium, MD",
	"Lutz, FL",
	"Lynbrook, NY",
	"Lyndhurst, NJ",
	"Mableton, GA",
	"Macon, GA",
	"Madison, WI",
	"Mahwah, NJ",
	"Maitland, FL",
	"Malden, MA",
	"Malibu, CA",
	"Malta, NY",
	"Malvern, PA",
	"Manasquan, NJ",
	"Manassas, VA",
	"Manchester, WA",
	"Manhasset, NY",
	"Manhattan Beach, CA",
	"Manhattan, NY",
	"Mankato, MN",
	"Mansfield, TX",
	"Maplewood, NJ",
	"Marblehead, MA",
	"Maricopa, AZ",
	"Marietta, GA",
	"Marina Del Rey, CA",
	"Marina, CA",
	"Marion, VA",
	"Marlborough, MA",
	"Maryland Heights, MO",
	"Maryville, MO",
	"Mason, OH",
	"Matawan, NJ",
	"Mayfield, KY",
	"Maywood, NJ",
	"Mazeppa, MN",
	"Mcdonough, GA",
	"Mchenry, IL",
	"Mckinney, TX",
	"Mclean, VA",
	"Mead, CO",
	"Mechanicsburg, PA",
	"Medfield, MA",
	"Medford, NJ",
	"Medford, OR",
	"Medway, MA",
	"Melbourne, FL",
	"Melville, NY",
	"Memphis, TN",
	"Menifee, CA",
	"Menlo Park, CA",
	"Menomonee Falls, WI",
	"Merchantville, NJ",
	"Merrick, NY",
	"Merrimack, NH",
	"Mesa, AZ",
	"Mesquite, TX",
	"Metuchen, NJ",
	"Miami Beach, FL",
	"Miami, FL",
	"Mid Town, NY",
	"Middleton, WI",
	"Middletown, NJ",
	"Midlothian, VA",
	"Milford, MI",
	"Mill Valley, CA",
	"Millbrae, CA",
	"Milpitas, CA",
	"Milwaukee, WI",
	"Mims, FL",
	"Minneapolis, MN",
	"Minnetonka, MN",
	"Miramar, FL",
	"Mission Viejo, CA",
	"Missouri City, TX",
	"Mobile, AL",
	"Mohegan Lake, NY",
	"Moline, IL",
	"Monmouth Junction, NJ",
	"Monroe Township, NJ",
	"Monroe, GA",
	"Montclair, CA",
	"Montclair, NJ",
	"Montebello, CA",
	"Monterey Park, CA",
	"Montrose, NY",
	"Montvale, NJ",
	"Monument, CO",
	"Moonachie, NJ",
	"Mooresville, NC",
	"Moorpark, CA",
	"Morgan Hill, CA",
	"Morgantown, WV",
	"Morganville, NJ",
	"Morris Plains, NJ",
	"Morristown, NJ",
	"Morristown, NJ",
	"Morristown, PA",
	"Morrisville, PA",
	"Morton Grove, IL",
	"Mount Airy, MD",
	"Mount Laurel, NJ",
	"Mount Pleasant, SC",
	"Mount Prospect, IL",
	"Mountain View, CA",
	"Mountainside, NJ",
	"Mountlake Terrace, WA",
	"Mountville, PA",
	"Moutain View, CA",
	"Mt Laurel, NJ",
	"Mukilteo, WA",
	"Mullica Hill, NJ",
	"Murray Hill, NJ",
	"Murrieta, CA",
	"Muskegon, MI",
	"Nanuet, NY",
	"Napa, CA",
	"Naperville, IL",
	"Naples, FL",
	"Nashua, NH",
	"Nashville, TN",
	"Natick, MA",
	"National City, CA",
	"Naugatuck, CT",
	"Navarre, FL",
	"Needham Heights, MA",
	"Neenah, WI",
	"Nevada City, CA",
	"New Bedford, MA",
	"New Brunswick, NJ",
	"New Haven, CT",
	"New Hyde Park, NY",
	"New London, PA",
	"New Orleans, LA",
	"New Paltz, NY",
	"New Port Richey, FL",
	"New Rochelle, NY",
	"New York, NY",
	"Newark, DE",
	"Newark, NJ",
	"Newberg, OR",
	"Newbury Park, CA",
	"Newington, NH",
	"Newport Beach, CA",
	"Newport, CA",
	"Newton Center, MA",
	"Newton Highlands, MA",
	"Newton, MA",
	"Norcross, GA",
	"Norfolk, VA",
	"North Arlington, NJ",
	"North Babylon, NY",
	"North Bend, WA",
	"North Billerica, MA",
	"North Brunswick, NJ",
	"North Easton, MA",
	"North Grosvenordale, CT",
	"North Hollywood, CA",
	"North Palm Beach, FL",
	"North Reading, MA",
	"North Richland Hills, TX",
	"Northborough, MA",
	"Northbrook, IL",
	"Northern, MI",
	"Northfield, MN",
	"Northford, CT",
	"Northport, NY",
	"Northridge, CA",
	"Northvale, NJ",
	"Norwalk, CT",
	"Norwell, MA",
	"Norwood, NJ",
	"Novi, MI",
	"Oak Brook, IL",
	"Oak Lawn, IL",
	"Oak Park, IL",
	"Oakland, CA",
	"Oakland, NJ",
	"Oakley, CA",
	"Oceanside, CA",
	"Ogden, UT",
	"Oklahoma City, OK",
	"Olathe, KS",
	"Old Bridge, NJ",
	"Oldsmar, FL",
	"Olympia, WA",
	"Omaha, NE",
	"Ontario, CA",
	"Oradell, NJ",
	"Orange County, CA",
	"Orem, UT",
	"Orinda, CA",
	"Orland Park, IL",
	"Orlando, FL",
	"Oshkosh, WI",
	"Overland Park, KS",
	"Owings Mills, MD",
	"Oxnard, CA",
	"Pacifica, CA",
	"Paia, HI",
	"Palm Beach Gardens, FL",
	"Palm Springs, CA",
	"Palmdale, CA",
	"Palmyra, PA",
	"Palo Alto, CA",
	"Palos Verdes Peninsula, CA",
	"Panama City, FL",
	"Paramount, CA",
	"Paramus, NJ",
	"Paris, TN",
	"Park City, UT",
	"Parlin, NJ",
	"Parsippany, NJ",
	"Pasadena, CA",
	"Payson, AZ",
	"Peachtree City, GA",
	"Pelham, NH",
	"Pennington, NJ",
	"Peoria, IL",
	"Pepperell, MA",
	"Petaluma, CA",
	"Pflugerville, TX",
	"Philadelphia, PA",
	"Phildelphia, PA",
	"Phillipsburg, NJ",
	"Phoenix, AZ",
	"Phoneix, AZ",
	"Piedmont, SC",
	"Pine Brook, NJ",
	"Pine River, MN",
	"Pineville, NC",
	"Piscataway, NJ",
	"Pismo Beach, CA",
	"Pittsburg, CA",
	"Pittsburgh, PA",
	"Plainsboro, NJ",
	"Plainview, NY",
	"Plano, TX",
	"Plantation, FL",
	"Plattsburgh, NY",
	"Playa Del Rey, CA",
	"Playa Vista, CA",
	"Pleasant Grove, UT",
	"Pleasant Hill, CA",
	"Pleasant Prairie, WI",
	"Pleasant Valley, NY",
	"Pleasanton, CA",
	"Pleasonton, CA",
	"Plymouth, MN",
	"Pompano Beach, FL",
	"Pooler, GA",
	"Port Orange, FL",
	"Port Washington, NY",
	"Port Washinton, NY",
	"Portage, IN",
	"Porter Ranch, CA",
	"Portland, OR",
	"Post Falls, ID",
	"Potomac, MD",
	"Poughkeepsie, NY",
	"Poway, CA",
	"Powder Springs, GA",
	"Prescott, AZ",
	"Princeton Junction, NJ",
	"Princeton, NJ",
	"Prospect Heights, IL",
	"Providence, RI",
	"Provo, UT",
	"Pullman, WA",
	"Purchase, NY",
	"Puyallup, WA",
	"Quantico, VA",
	"Queen Creek, AZ",
	"Quincy, MA",
	"Raleigh, NC",
	"Ramona, CA",
	"Rancho Bernardo, CA",
	"Rancho Cordova, CA",
	"Rancho Cucamonga, CA",
	"Rancho Santa Fe, CA",
	"Rapid City, SD",
	"Raritan, NJ",
	"Reading, MA",
	"Reading, PA",
	"Red Bank, NJ",
	"Redding, CA",
	"Redlands, CA",
	"Redmond, WA",
	"Redondo Beach, CA",
	"Redwood City, CA",
	"Redwood Shores, CA",
	"Reno, NV",
	"Renton, WA",
	"Reseda, CA",
	"Reston, VA",
	"Rexburg, ID",
	"Reynoldsburg, OH",
	"Rialto, CA",
	"Richardson, TX",
	"Richfield, MN",
	"Richland, WA",
	"Richmond, IN",
	"Richmond, VA",
	"Ridgefield Park, NJ",
	"Ridgefield, WA",
	"Riverside, CT",
	"Riverton, UT",
	"Riverwoods, IL",
	"Roanoke, TX",
	"Rochester, NY",
	"Rock Hill, SC",
	"Rockville, MD",
	"Rogers, AR",
	"Rohnert Park, CA",
	"Romeoville, IL",
	"Ronkonkoma, NY",
	"Roseland, NJ",
	"Roselle, IL",
	"Rosemead, CA",
	"Rosenberg, TX",
	"Roseville, MN",
	"Roslyn Heights, NY",
	"Roslyn, NY",
	"Rosslyn, VA",
	"Roswell, GA",
	"Round Rock, TX",
	"Rowland Heights, CA",
	"Rowlett, TX",
	"Rowley, MA",
	"Ruston, LA",
	"Rutherford, NJ",
	"Rye, NY",
	"Sacramento, CA",
	"Sahuarita, AZ",
	"Saint Augustine, FL",
	"Saint Cloud, MN",
	"Saint George, UT",
	"Saint Johns, FL",
	"Saint Joseph, MO",
	"Saint Louis, MO",
	"Saint Michael, MN",
	"Saint Paul, MN",
	"Saint Petersburg, FL",
	"Salem, MA",
	"Salt Lake City, UT",
	"Sammamish, WA",
	"San Angelo, TX",
	"San Antonio, TX",
	"San Bernardino, CA",
	"San Bruno, CA",
	"San Carlos City, CA",
	"San Carlos, CA",
	"San Deigo, CA",
	"San Diego, CA",
	"San Dimas, CA",
	"San Francisco, CA",
	"San Gabriel, CA",
	"San Jacinto, CA",
	"San Jose, CA",
	"San Juan Capistrano, CA",
	"San Leandro, CA",
	"San Luis Obispo, CA",
	"San Marcos, CA",
	"San Mateo, CA",
	"San Pablo, CA",
	"San Rafael, CA",
	"San Ramon, CA",
	"Sandy Hook, CT",
	"Sanford, ME",
	"Santa Ana, CA",
	"Santa Barbara, CA",
	"Santa Clara, CA",
	"Santa Clarita, Ca",
	"Santa Cruz, CA",
	"Santa Fe, NM",
	"Santa Monica, CA",
	"Santee, CA",
	"Sarasota, FL",
	"Saratoga Springs, NY",
	"Saratoga, CA",
	"Sault Sainte Marie, MI",
	"Sausalito, CA",
	"Savage, MN",
	"Savannah, GA",
	"Scarsdale, NY",
	"Schaumburg, IL",
	"Scotts Valley, CA",
	"Scottsdale, AZ",
	"Scottsdate, AZ",
	"Scranton, PA",
	"Seaside, CA",
	"Seattle, WA",
	"Secaucus, NJ",
	"Secausus, NJ",
	"Selden, NY",
	"Severna Park, MD",
	"Sharon, MA",
	"Shelton, CT",
	"Sherman Oaks, CA",
	"Short Hills, NJ",
	"Shreve, OH",
	"Shrewsbury, MA",
	"Sierra Madre, CA",
	"Sierra Vista, AZ",
	"Silicon Valley, CA",
	"Silver Spring, MD",
	"Simi Valley, CA",
	"Sioux Falls, SD",
	"Skillman, NJ",
	"Skokie, IL",
	"Smithfield, RI",
	"Smyrna, GA",
	"Snohomish, WA",
	"Snoqualmie, WA",
	"Solana Beach, CA",
	"Solon, OH",
	"Somerset, NJ",
	"Somerville, MA",
	"Somerville, NJ",
	"Sound Beach, NY",
	"South Bound Brook, NJ",
	"South Grafton, MA",
	"South Jordan, UT",
	"South Ozone Park, NY",
	"South Pasadena, CA",
	"Southboro, MA",
	"Southborough, MA",
	"Southbury, CT",
	"Southfield, MI",
	"Southlake, TX",
	"Spokane, WA",
	"Spring Valley, NY",
	"Spring, TX",
	"Springfield, VA",
	"Springville, UT",
	"Stafford, VA",
	"Stamford, CT",
	"Stanford, CA",
	"Starkville, MS",
	"Staten Island, NY",
	"Sterling, VA",
	"Stillwater, OK",
	"Stonington, CT",
	"Stony Brook, NY",
	"Stony Point, NY",
	"Storrs, CT",
	"Strongsville, OH",
	"Studio City, CA",
	"Sudbury, MA",
	"Sugar Land, TX",
	"Sun City, AZ",
	"Sun Prairie, WI",
	"Sunnyside, NY",
	"Sunnyvale, CA",
	"Superior, CO",
	"Surprise, AZ",
	"Suwanee, GA",
	"Swanton, VT",
	"Sykesville, MD",
	"Sylmar, CA",
	"Sylva, NC",
	"Syracuse, NY",
	"Tacoma, WA",
	"Takoma Park, MD",
	"Tallahassee, FL",
	"Tampa, FL",
	"Tappan, NY",
	"Tarrytown, NY",
	"Taylor, PA",
	"Taylorsville, UT",
	"Teaneck, NJ",
	"Temecula, CA",
	"Tempe, AZ",
	"Temple City, CA",
	"Temple Terrace, FL",
	"Tenafly, NJ",
	"Thousand Oaks, CA",
	"Tinton Fall, NJ",
	"Tinton Falls, NJ",
	"Titusville, NJ",
	"Toledo, OH",
	"Tonawanda, NY",
	"Topeka, KS",
	"Torrance, CA",
	"Torrence, CA",
	"Totowa, NJ",
	"Trabuco Canyon, CA",
	"Tracy, CA",
	"Traverse City, MI",
	"Trenton, NJ",
	"Troy, AL",
	"Trumbull, CT",
	"Tuckahoe, NY",
	"Tucson, AZ",
	"Tulsa, OK",
	"Tustin, CA",
	"Tysons Corner, VA",
	"Union Bridge, MD",
	"Union City, CA",
	"Union City, NJ",
	"Union, MO",
	"Union, NJ",
	"Uniondale, NY",
	"Universal City, CA",
	"Upland, IN",
	"Upper Marlboro, MD",
	"Upton, MA",
	"Vacaville, CA",
	"Vail, AZ",
	"Valencia, CA",
	"Valley Stream, NY",
	"Valley Village, CA",
	"Valley, NE",
	"Van Nuys, CA",
	"Vancouver, WA",
	"Venice, CA",
	"Ventura, CA",
	"Vernon Hills, IL",
	"Vernon, CA",
	"Vestal, NY",
	"Vicksburg, MS",
	"Victor, NY",
	"Vienna, VA",
	"Virginia Beach, VA",
	"Vista, CA",
	"Voorhees, NJ",
	"Waddell, AZ",
	"Wake Forest, NC",
	"Wakefield, MA",
	"Walkersville, MD",
	"Wallingford, PA",
	"Walnut Creek, CA",
	"Walnut, CA",
	"Waltham, MA",
	"Warren, NJ",
	"Warrenville, IL",
	"Washington, DC",
	"Watertown, MA",
	"Waukesha, WI",
	"Waxhaw, NC",
	"Wayne, NJ",
	"Wayne, PA",
	"Wayzata, MN",
	"Weehawken, NJ",
	"Wellesley Hills, MA",
	"Wesley Chapel, FL",
	"West Chester, PA",
	"West Columbia, SC",
	"West Covina, CA",
	"West Des Moines, IA",
	"West Haven, CT",
	"West Hills, CA",
	"West Hollywood, CA",
	"West Jordan, UT",
	"West Lafayette, IN",
	"West Liberty, KY",
	"West Los Angeles, CA",
	"West New York, NJ",
	"West Newbury, MA",
	"West Newton, MA",
	"West Orange, NJ",
	"West Palm Beach, FL",
	"Westborough, MA",
	"Westcliffe, CO",
	"Westfield, MA",
	"Westford, MA",
	"Westlake Village, CA",
	"Westlake, TX",
	"Westminster, CO",
	"Weston, FL",
	"Westport, CT",
	"Westwood, NJ",
	"Wharton, NJ",
	"White Plains, NY",
	"Whitestone, NY",
	"Wildomar, CA",
	"Willow Grove, PA",
	"Wilmington, NC",
	"Wilton, CT",
	"Wimberley, TX",
	"Windermere, FL",
	"Windsor Locks, CT",
	"Windsor, CO",
	"Winnetka, CA",
	"Winston Salem, NC",
	"Winter Park, FL",
	"Woburn, MA",
	"Woodbridge, NJ",
	"Woodbridge, VA",
	"Woodbury, NY",
	"Woodcliff Lake, NJ",
	"Woodcliffe Lake, NJ",
	"Woodinville, WA",
	"Woodland Hills, CA",
	"Woodstock, GA",
	"Worcester, MA",
	"Wylie, TX",
	"Yonkers, NY",
	"Yorba Linda, CA",
	"Yorktown Heights, NY",
	"Yorktown, NY",
	"Yuba City, CA",
}

// this table is not used, only serves ex_main:
var nonStates = []string{
	// leave 1 active, lest regexp fail...
	"Washington D.C.", // incl. "Washington D.C. Metro Area"
	//	"Washington, District Of Columbia",
	//	"Washington, DC",
	//	"New York City", // incl. "Greater New York City Area"
	//	"New York, New York",
	//	"New York, NY",
	//	"NYC",
	//	"Kansas City",
	//	"Missouri City",
	//	"Nevada City",
	//	"Oklahoma City",
	//	"Iowa City",
	//	"Colorado Springs",
	//	"Idaho Falls",
	//	"DC Metro",
	//	"Maryland Heights",
	//	"Manhattan Beach",
	//	"Virginia Beach",
	//	"Port Washington",
	//	"Indianapolis",
	//	"Hawaiian Islands",
	//	"SF Bay",
	//
	//	// many of those I could use, just must solve 2 of these:
	//
	//	"Hawaiian Islands",
	//
	//	"Washington D.C.", // incl. "Washington D.C. Metro Area"
	//	"Washington, District Of Columbia",
	//	"Washington, DC",
	//	"DC Metro",
	//
	//	"SF Bay",
	//
	//	"New York City", // incl. "Greater New York City Area"
	//	"New York, New York",
	//	"New York, NY",
	//	"NYC",
}

var states = []string{
	"AK,Alaska",
	"AL,Alabama",
	"AR,Arkansas",
	"AZ,Arizona",
	"CA,California",
	"CO,Colorado",
	"CT,Connecticut",
	"DC,District of Columbia",
	"DE,Delaware",
	"FL,Florida",
	"GA,Georgia",
	"HI,Hawaii",
	"IA,Iowa",
	"ID,Idaho",
	"IL,Illinois",
	"IN,Indiana",
	"KS,Kansas",
	"KY,Kentucky",
	"LA,Louisiana",
	"MA,Massachusetts",
	"MD,Maryland",
	"ME,Maine",
	"MI,Michigan",
	"MN,Minnesota",
	"MO,Missouri",
	"MS,Mississippi",
	"MT,Montana",
	"NC,North Carolina",
	"ND,North Dakota",
	"NE,Nebraska",
	"NH,New Hampshire",
	"NJ,New Jersey",
	"NM,New Mexico",
	"NV,Nevada",
	"NY,New York",
	"OH,Ohio",
	"OK,Oklahoma",
	"OR,Oregon",
	"PA,Pennsylvania",
	"RI,Rhode Island",
	"SC,South Carolina",
	"SD,South Dakota",
	"TN,Tennessee",
	"TX,Texas",
	"UT,Utah",
	"VA,Virginia",
	"VT,Vermont",
	"WA,Washington",
	"WI,Wisconsin",
	"WV,West Virginia",
	"WY,Wyoming",
}

// Foreign countries match streets, cities, states, a delicate matter.
var countries = []string{
	// Manually added:
	"Russian Federation",
	// Foreign Country list found on the web:
	"China",
	// Hey, The U.S. is not foreign! -- "United States",
	"European Union",
	"India(\\b|$)", // HID Avon, Indiana and others. Fix with \b, but also $
	"ASEAN",
	"Japan",
	"Germany",
	"EAEU",
	"Russia",
	"Indonesia",
	"United Kingdom",
	"Brazil",
	"France",
	"Turkey",
	"Italy",
	// omit, hides US state New Mexico -- "Mexico",
	"South Korea",
	"Canada",
	"Spain",
	"Saudi Arabia",
	"Australia",
	"Taiwan",
	"Poland",
	"Iran",
	"Egypt",
	"Thailand",
	"Pakistan",
	"Vietnam",
	"Nigeria",
	"Netherlands",
	"Argentina",
	"Philippines",
	"Bangladesh",
	"Malaysia",
	"Colombia",
	"South Africa",
	"United Arab Emirates",
	"Switzerland",
	"Belgium",
	"Romania",
	"Singapore",
	"Sweden",
	"Ireland",
	"Kazakhstan",
	"Ukraine",
	"Algeria",
	"Austria",
	"Chile",
	"Hong Kong",
	"Peru",
	"Iraq",
	"Czech Republic",
	"Israel",
	"Norway",
	"Portugal",
	"Denmark",
	"Hungary",
	"Greece",
	"Ethiopia",
	"Sri Lanka",
	"Morocco",
	"Uzbekistan",
	"Finland",
	"Kenya",
	"Qatar",
	"New Zealand",
	"Myanmar",
	"Dominican Republic",
	"Kuwait",
	"Angola",
	"Ecuador",
	"Ghana",
	"Slovakia",
	"Sudan",
	"Tanzania",
	"Belarus",
	"Bulgaria",
	"Guatemala",
	"Ivory Coast",
	"Azerbaijan",
	"Oman",
	"Serbia",
	"Venezuela",
	// hides a city in FL -- "Panama",
	"Tunisia",
	"Croatia",
	"Nepal",
	"Cuba",
	"Puerto Rico",
	"Lithuania",
	"Uganda",
	"Costa Rica",
	"DR Congo",
	"Libya",
	"Cameroon",
	// Omit, hides 2 cities in UT -- "Jordan",
	"Bolivia",
	"Turkmenistan",
	"Paraguay",
	"Slovenia",
	"Uruguay",
	"Luxembourg",
	"Cambodia",
	"Bahrain",
	// hides a city in NJ -- "Lebanon",
	"Afghanistan",
	"Zambia",
	"Senegal",
	"Latvia",
	"Honduras",
	"El Salvador",
	// omit -- "Georgia", // lest this country hide USA State
	"Laos",
	"Yemen",
	"Bosnia and Herzegovina",
	"Macau",
	"Estonia",
	"Burkina Faso",
	// Hides Malibu CA -- "Mali",
	"Benin",
	"Madagascar",
	"Syria",
	"Albania",
	"Mozambique",
	"Botswana",
	"Armenia",
	"Nicaragua",
	"Mongolia",
	"Tajikistan",
	"Guinea",
	"Cyprus",
	"Moldova",
	"Trinidad and Tobago",
	"North Macedonia",
	"North Korea",
	"Zimbabwe",
	"Papua New Guinea",
	"Gabon",
	"Haiti",
	"Kyrgyzstan",
	"Niger",
	"Rwanda",
	"Malawi",
	"Palestine",
	"Brunei",
	"Jamaica",
	"Mauritius",
	"Guyana",
	"Mauritania",
	"Chad",
	"Equatorial Guinea",
	// hides a city in NY -- "Malta",
	"Namibia",
	"Kosovo",
	"Iceland",
	"Togo",
	"Congo",
	"Somalia",
	"Sierra Leone",
	"Bahamas",
	"Montenegro",
	"South Sudan",
	"Fiji",
	"Eswatini",
	"Maldives",
	"New Caledonia",
	"Burundi",
	"Suriname",
	"Bhutan",
	"Liberia",
	"Eritrea",
	"Monaco",
	"Isle of Man",
	"Gambia",
	"Djibouti",
	"Lesotho",
	"Guam",
	// Omit, hides New Jersey -- "Jersey",
	"Central African Republic",
	"French Polynesia",
	"Bermuda",
	"Guinea-Bissau",
	"Andorra",
	"Barbados",
	"Liechtenstein",
	"Cayman Islands",
	"East Timor",
	"Aruba",
	"Cape Verde",
	"Curaçao",
	"U.S. Virgin Islands",
	"Seychelles",
	"Guernsey",
	"Comoros",
	"Belize",
	"Saint Lucia",
	"Greenland",
	"San Marino",
	"Antigua and Barbuda",
	"Gibraltar",
	"Grenada",
	"Faroe Islands",
	"Saint Vincent and the Grenadines",
	"Solomon Islands",
	"Saint Kitts and Nevis",
	"Sint Maarten",
	"Northern Mariana Islands",
	"Samoa",
	"Dominica",
	"São Tomé and Príncipe",
	"Vanuatu",
	"Turks and Caicos Islands",
	"Tonga",
	"American Samoa",
	"Saint Martin",
	"British Virgin Islands",
	"Micronesia",
	"Cook Islands",
	"Kiribati",
	"Palau",
	"Saint Pierre and Miquelon",
	"Marshall Islands",
	"Falkland Islands",
	"Anguilla",
	"Montserrat",
	"Nauru",
	"Wallis and Futuna",
	"Tuvalu",
	"Saint Helena, Ascension and Tristan da Cunha",
	"Niue",
	"Tokelau",
}

// This is dependent on the data I scraped.
// A few ringers crept in, I had to remove.
// A few others I must manually add herein.

// This list is not used, only serves ex_main 1 and/or 2.

// var QualityCities = []string{
var QualityCities = S2S{
	//	// manual additions of cities resembling states
	// leave 1 active, lest regexp fail...
	"Kansas City",
	//	"Missouri City",
	//	"Nevada City",
	//	"Oklahoma City",
	//	"Iowa City",
	//	"Colorado Springs",
	//	"Idaho Falls",
	//	"Maryland Heights",
	//	"Manhattan Beach",
	//	"Virginia Beach",
	//	"Port Washington",
	//	"Indianapolis",
	//	"Chesterbrook", // wrongly rid as a school, but is also a city
	//	// cpu-generated, vetted part
	//	"Phoenix",
	//	"San Diego",
	//	"San Jose",
	//	"San Francisco",
	//	"Sunnyvale",
	//	"Austin",
	//	"Mountain View",
	//	"Santa Clara",
	//	"Los Angeles",
	//	"Tempe",
	//	"Palo Alto",
	//	"Scottsdale",
	//	"Chicago",
	//	"Seattle",
	//	"Raleigh",
	//	"Irvine",
	//	"Orange County",
	//	"Charlotte",
	//	"Durham",
	//	"Chandler",
	//	"Houston",
	//	"Portland",
	//	"Redwood City",
	//	"Boston",
	//	"Cupertino",
	//	"Fremont",
	//	"Atlanta",
	//	"Jersey City",
	//	"Tampa",
	//	"San Ramon",
	//	"Carlsbad",
	//	"Pleasanton",
	//	"Dallas",
	//	"San Mateo",
	//	"Santa Monica",
	//	"Tucson",
	//	"Saint Petersburg",
	//	"Redmond",
	//	"Bellevue",
	//	"Menlo Park",
	//	"Orlando",
	//	"Brooklyn",
	//	"Plano",
	//	"Alpharetta",
	//	"Columbus",
	//	"Baltimore",
	//	"Irving",
	//	"San Antonio",
	//	"Herndon",
	//	"Milpitas",
	//	"Philadelphia",
	//	"Reston",
	//	"Oakland",
	//	"Richardson",
	//	"Arlington",
	//	"Pasadena",
	//	"Cambridge",
	//	"Rochester",
	//	"Glendale",
	//	"Redwood Shores",
	//	"Denver",
	//	"Boulder",
	//	"Nashville",
	//	"La Jolla",
	//	"Cleveland",
	//	"Minneapolis",
	//	"Richmond",
	//	"Gilbert",
	//	"Princeton",
	//	"Akron",
	//	"Los Gatos",
	//	"Edison",
	//	"Cincinnati",
	//	"San Marcos",
	//	"Mesa",
	//	"Sacramento",
	//	"El Segundo",
	//	"Pittsburgh",
	//	"Burbank",
	//	"Louisville",
	//	"Las Vegas",
	//	"Berkeley",
	//	// UNUNIQUE -- "Mclean",
	//	"Foster City",
	//	"Cary",
	//	"Emeryville",
	//	"Woodland Hills",
	//	"McLean",
	//	"Hartford",
	//	"Bentonville",
	//	"Culver City",
	//	"Miami",
	//	"Albany",
	//	"Fairfax",
	//	"Encinitas",
	//	"Union City",
	//	"Columbia",
	//	"Poway",
	//	"Torrance",
	//	"Jacksonville",
	//	"Campbell",
	//	"Newark",
	//	"Madison",
	//	"Dublin", // matches County Dublin, Ireland -- NP, rid foreign first
	//	"Costa Mesa",
	//	"Stamford",
	//	"Frisco",
	//	"Littleton",
	//	"Santa Barbara",
	//	"Provo",
	//	"Wilmington",
	//	"Oceanside",
	//	"Venice",
	//	"Belmont",
	//	"Hoboken",
	//	"Solana Beach",
	//	"Waltham",
	//	"Piscataway",
	//	"Bloomington",
	//	"Los Altos",
	//	"Fort Collins",
	//	"Schaumburg",
	//	"Thousand Oaks",
	//	"Kirkland",
	//	"Ashburn",
	//	"Fort Lauderdale",
	//	"Hollywood",
	//	"Burlington",
	//	"Providence",
	//	"Middletown",
	//	"Rockville",
	//	"Peoria",
	//	"Saint Louis",
	//	"Buffalo",
	//	"Warren",
	//	"Hopkinton",
	//	"Greenville",
	//	"Ann Arbor",
	//	"Salt Lake City",
	//	"Charleston",
	//	"Basking Ridge",
	//	"Long Beach",
	//	"Framingham",
	//	"Englewood",
	//	"Plainsboro",
	//	"Grand Rapids",
	//	"Gainesville",
	//	"Vienna",
	//	"Aliso Viejo",
	//	"White Plains",
	//	"Stony Brook",
	//	"Lexington",
	//	"Escondido",
	//	// ringer -- "Southern",
	//	"West Palm Beach",
	//	"San Bruno",
	//	"Omaha",
	//	"Walnut Creek",
	//	"Santa Cruz",
	//	"Boca Raton",
	//	"Harrisburg",
	//	"Sterling",
	//	"San Carlos",
	//	"Fort Worth",
	//	"Vista",
	//	"Somerset",
	//	"Parsippany",
	//	"Fayetteville",
	//	"Bridgewater",
	//	"College Park",
	//	"Beaverton",
	//	"Harrison",
	//	"Bridgeport",
	//	"Chula Vista",
	//	"Livermore",
	//	"Lehi",
	//	"Hillsboro",
	//	"Hayward",
	//	"Gaithersburg",
	//	"Frederick",
	//	"Santee",
	//	"College Station",
	//	"Broomfield",
	//	"Marina Del Rey",
	//	"Sarasota",
	//	"Beverly Hills",
	//	"Huntington Beach",
	//	"Birmingham",
	//	"Round Rock",
	//	"Redlands",
	//	"Tulsa",
	//	"Temecula",
	//	"Lake Forest",
	//	"Riverside",
	//	"Naperville",
	//	"Dayton",
	//	"Draper",
	//	"Woodbridge",
	//	"Eden Prairie",
	//	"Niagara",
	//	"Islandia",
	//	"Silver Spring",
	//	"Cranbury",
	//	"Hicksville",
	//	"Duluth",
	//	"Springfield",
	//	"Saratoga",
	//	"Bethesda",
	//	"Anaheim",
	//	"Greensboro",
	//	"Winston",
	//	"Dover",
	//	"Salem",
	//	// ringer -- "South San Francisco",
	//	"Nashua",
	//	"Memphis",
	//	"Leesburg",
	//	"Calabasas",
	//	// ringer -- "Remote",
	//	"Manhattan",
	//	"Wilton",
	//	"Norfolk",
	//	"Morrisville",
	//	"Charlottesville",
	//	"Hoffman Estates",
	//	"Danville",
	//	"Fullerton",
	//	"Weehawken",
	//	"Simi Valley",
	//	"Alexandria",
	//	"Boise",
	//	"Concord",
	//	// ringer -- "San Francisco Bay",
	//	"Syracuse",
	//	"Poughkeepsie",
	//	"Chantilly",
	//	"Daly City",
	//	"Worcester",
	//	"Evanston",
	//	"Norwalk",
	//	"Huntsville",
	//	"El Cajon",
	//	"Roseland",
	//	"Windsor",
	//	"Richfield",
	//	"Armonk",
	//	"Burlingame",
	//	"Allen",
	//	"Bryan",
	//	"Lansing",
	//	"Tarrytown",
	//	"Richland",
	//	"Melbourne",
	//	"West Chester",
	//	"Pennington",
	//	"Allentown",
	//	// ringer -- "Research Triangle Park",
	//	"East Brunswick",
	//	"Bronx",
	//	"Westlake Village",
	//	"San Leandro",
	//	"Astoria",
	//	"Encino",
	//	"Morgan Hill",
	//	"Bedford",
	//	"Marietta",
	//	"Annapolis",
	//	"Yorktown Heights",
	//	"Stanford",
	//	"Milford",
	//	"Corona",
	//	"Playa Vista",
	//	"Bedminster",
	//	"Suwanee",
	//	"Lawrence",
	//	"Toledo",
	//	"Manchester",
	//	"Overland Park",
	//	"Jackson",
	//	"Medford",
	//	"Normal",
	//	"Fairfield",
	//	"West Hollywood",
	//	"Greenwich",
	//	"Knoxville",
	//	"Westminster",
	//	"Aurora",
	//	"Albuquerque",
	//	"Roswell",
	//	"Westborough",
	//	"Melville",
	//	"El Paso",
	//	"Sammamish",
	//	"Mahwah",
	//	"New Brunswick",
	//	"Savannah",
	//	"Falls Church",
	//	"Buffalo Grove",
	//	"Tacoma",
	//	"Miami Beach",
	//	"Iselin",
	//	"Mason",
	//	"Wayne",
	//	"Tenafly",
	//	"Dearborn",
	//	"New Haven",
	//	"Lutz",
	//	"Agoura Hills",
	//	"Greenwood Village",
	//	"Eugene",
	//	"Ames",
	//	"Lafayette",
	//	"Rancho Cucamonga",
	//	"Oak Park",
	//	"Somerville",
	//	"Manassas",
	//	"Amherst",
	//	"Davis",
	//	"Cave Creek",
	//	"Trenton",
	//	"Millbrae",
	//	"Sugar Land",
	//	"East Lansing",
	//	"Orem",
	//	"Laveen",
	//	"Maitland",
	//	"Bothell",
	//	"Pasco",
	//	"Berkeley Heights",
	//	"Mooresville",
	//	"Fort Smith",
	//	"Northridge",
	//	"Sherman Oaks",
	//	"Potomac",
	//	"Fredericksburg",
	//	"Kennewick",
	//	"Union",
	//	"Monroe Township",
	//	"Westport",
	//	"Scotts Valley",
	//	"Ellicott City",
	//	"Champaign",
	//	"Ventura",
	//	"Mill Valley",
	//	"Trumbull",
	//	"Shelton",
	//	"Scranton",
	//	"Glen Allen",
	//	"La Mesa",
	//	"Henderson",
	//	"Renton",
	//	"Baton Rouge",
	//	"Little Elm",
	//	"North Hollywood",
	//	"Gilroy",
	//	"Detroit",
	//	// UNUNIQUE -- "Mountain view",
	//	"Spring Valley",
	//	"Brea", // hmmm... matched Panera Bread. should add \b, retry.
	//	"Yonkers",
	//	"Fort Mill",
	//	"Little Rock",
	//	"Folsom",
	//	"Carrollton",
	//	"Camp Hill",
	//	"Quincy",
	//	"Northbrook",
	//	"Holmdel",
	//	"Germantown",
	//	"Franklin Lakes",
	//	"Oak Brook",
	//	"Hillsborough",
	//	"Morristown",
	//	"Claremont",
	//	"Fords",
	//	"Cumming",
	//	"Cerritos",
	//	"Alameda",
	//	"La Palma",
	//	"Clearwater",
	//	"Rye",
	//	"Rexburg",
	//	"Marlborough",
	//	"Newport Beach",
	//	"Des Moines",
	//	"Napa",
	//	"Fort Walton Beach",
	//	"Lincoln",
	//	"Santa Ana",
	//	"El Cerrito",
	//	"Reading",
	//	"Apex",
	//	"Indian Trail",
	//	"Beaumont",
	//	"Lincolnshire",
	//	"Teaneck",
	//	"Urbana",
	//	"Miramar",
	//	"Colleyville",
	//	// ringer -- "West",
	//	"Jupiter",
	//	"Oshkosh",
	//	"Bohemia",
	//	"Van Nuys",
	//	"Lake Mary",
	//	"Lowell",
	//	"Macon",
	//	"Torrence",
	//	"Dekalb",
	//	"Long Island City",
	//	"West Orange",
	//	"Brisbane",
	//	"Pullman",
	//	"Kendall Park",
	//	"Camarillo",
	//	"Utica",
	//	"Collierville",
	//	// ringer -- "RTP",
	//	"Trabuco Canyon",
	//	"Watertown",
	//	"Barrington",
	//	"Saratoga Springs",
	//	"North Brunswick",
	//	"Murrieta",
	//	"Horsham",
	//	"Danbury",
	//	"Auburn",
	//	"Skokie",
	//	"Princeton Junction",
	//	"Georgetown",
	//	"Keller",
	//	"Conshohocken",
	//	"Brighton",
	//	"Hawthorne",
	//	"Issaquah",
	//	"Davenport",
	//	"Lubbock",
	//	"Newberg",
	//	"Westford",
	//	"Palmdale",
	//	"Valley Stream",
	//	"Canoga Park",
	//	"Sierra Vista",
	//	"Winter Park",
	//	"San Rafael",
	//	"Everett",
	//	"Pleasant Grove",
	//	"Flemington",
	//	"Alhambra",
	//	"Ridgefield",
	//	"Savage",
	//	"Universal City",
	//	"Short Hills",
	//	"Springville",
	//	"Saint Paul",
	//	"Tallahassee",
	//	"Canton",
	//	"Severna Park",
	//	"Sausalito",
	//	"Needham Heights",
	//	"Chelmsford",
	//	"Golden",
	//	"Cherry Hill",
	//	"Independence",
	//	"Centreville",
	//	"Cardiff By The Sea",
	//	"Addison",
	//	"Monmouth Junction",
	//	"Hightstown",
	//	"Queen Creek",
	//	"Fountain Valley",
	//	"Clifton",
	//	"Natick",
	//	"Covington",
	//	"Chambersburg",
	//	"Harrisonburg",
	//	"Montvale",
	//	"Laurel",
	//	"Dulles",
	//	"Binghamton",
	//	"Longmont",
	//	"Norcross",
	//	// UNUNIQUE -- "Foster city",
	//	"Smyrna",
	//	"Bensenville",
	//	"Englewood Cliffs",
	//	"Hamden",
	//	"Naples",
	//	"Great Falls",
	//	"Ashland",
	//	"Burnsville",
	//	"Hammond",
	//	"Newington",
	//	"Mckinney",
	//	"West Covina",
	//	"Delray Beach",
	//	"Stafford",
	//	"Upper Marlboro",
	//	"San Luis Obispo",
	//	"Palm Springs",
	//	"Jamestown",
	//	"Lawrenceburg",
	//	"Fairview",
	//	// UNUNIQUE -- "Redwood city",
	//	"Southboro",
	//	"Ruston",
	//	"Bloomfield",
	//	"Charlestown",
	//	"Sanford",
	//	"Reseda",
	//	"Saint George",
	//	"Bowie",
	//	"Wakefield",
	//	"Lake Success",
	//	// UNUNIQUE -- "Redwood shores",
	//	"Lakeside",
	//	"Castle Rock",
	//	"Moutain View",
	//	"Ramona",
	//	"Saint Cloud",
	//	"Adelphi",
	//	"Pine Brook",
	//	"Lake Oswego",
	//	"Victor",
	//	"Novi",
	//	"Chico",
	//	"Malden",
	//	"Raritan",
	//	"Huntersville",
	//	"Saint Augustine",
	//	"Lodi",
	//	"Fulton",
	//	"Riverwoods",
	//	"Southlake",
	//	"Glenview",
	//	"Panama City",
	//	"Mansfield",
	//	"Clayton",
	//	"Ontario",
	//	"Kissimmee",
	//	"Roselle",
	//	"Foothill Ranch",
	//	"Wellesley Hills",
	//	"Silicon Valley",
	//	"Malvern",
	//	"Lake Grove",
	//	"Bonsall",
	//	"Sandy Hook",
	//	"Kent", // This city, Kent, WA, matches state Kentucky. But \b may fix.
	//	"Destin",
	//	"La Mirada",
	//	"San Bernardino",
	//	"Titusville",
	//	"Woodstock",
	//	"Kenmore",
	//	"Lawrenceville",
	//	"Athens",
	//	"Johns Creek",
	//	"Reno",
	//	"Tustin",
	//	"Goodyear",
	//	"Haymarket",
	//	"Conroe",
	//	"Winnetka",
	//	"Roseville",
	//	"Westlake",
	//	"Woodcliff Lake",
	//	"San Dimas",
	//	"Carle Place",
	//	"Bolingbrook",
	//	"Moorpark",
	//	"Carmel",
	//	"King of Prussia",
	//	"Elmont",
	//	"Blacksburg",
	//	"Elkridge",
	//	"Huntington Station",
	//	"Ladera Ranch",
	//	"Saint Johns",
	//	"Merrimack",
	//	"Pleasant Hill",
	//	"Morganville",
	//	"Studio City",
	//	"Kearny",
	//	"Annandale",
	//	"Coppell",
	//	"Palm Beach Gardens",
	//	"Sahuarita",
	//	"Andover",
	//	"Rosemead",
	//	"Waxhaw",
	//	"Coralville",
	//	"Chapel Hill",
	//	"Southfield",
	//	"Valencia",
	//	"Acworth",
	//	"Middleton",
	//	"Hudson",
	//	"Dedham",
	//	"El Monte",
	//	"Fargo",
	//	"Fort Huachuca",
	//	"Farmington",
	//	"Buford",
	//	"Wildomar",
	//	"Augusta",
	//	"Buckeye",
	//	"Pflugerville",
	//	"Cliffside Park",
	//	"Kennebunk",
	//	"San Jacinto",
	//	"Carpinteria",
	//	"American Fork",
	//	// ringer -- "Phoenix y alrededores",
	//	"Bainbridge Island",
	//	"Briarcliff Manor",
	//	"North Grosvenordale",
	//	"Montclair",
	//	"Arvada",
	//	"Fishers",
	//	"Lewis Center",
	//	"Morgantown",
	//	"San Angelo",
	//	"Pacifica",
	//	"Rancho Cordova",
	//	"Newbury Park",
	//	"Tysons Corner",
	//	"San Juan Capistrano",
	//	"Coraopolis",
	//	"Windermere",
	//	"Metuchen",
	//	"Vancouver",
	//	"Bonita",
	//	"East Stroudsburg",
	//	"Pineville",
	//	"Florham Park",
	//	"Orinda",
	//	"Forest Hills",
	//	"Upton",
	//	"Logan",
	//	"Purchase",
	//	"High Point",
	//	"Milwaukee",
	//	"Midlothian",
	//	// UNUNIQUE -- "King Of Prussia",
	//	"New Hyde Park",
	//	"Redondo Beach",
	//	"Petaluma",
	//	"Northvale",
	//	"Snoqualmie",
	//	"Center Valley",
	//	"Del Mar",
	//	"West Lafayette",
	//	"Duarte",
	//	"Long Island",
	//	"Red Bank",
	//	"Pompano Beach",
	//	"Elk Grove",
	//	"Half Moon Bay",
	//	"La Crescenta",
	//	"Wake Forest",
	//	"Plymouth",
	//	"Arlington Heights",
	//	"Bozeman",
	//	"Topeka",
	//	"Olympia",
	//	"Federal Way",
	//	"Denton",
	//	"Mountlake Terrace",
	//	"Sharon",
	//	"Garden Grove",
	//	"Broken Arrow",
	//	"Horseshoe Bay",
	//	"Paramus",
	//	"Vacaville",
	//	"Cedar Rapids",
	//	"Benicia",
	//	"Hercules",
	//	"West Haven",
	//	"Staten Island",
	//	"Granite Bay",
	//	"Hingham",
	//	"Tinton Falls",
	//	"Sioux Falls",
	//	"Highlands Ranch",
	//	"Ronkonkoma",
	//	"Hempstead",
	//	"Holly Springs",
	//	"Great Barrington",
	//	"Cape Girardeau",
	//	"Green Bay",
	//	"Boonton",
	//	"Des Plaines",
	//	"Lindon",
	//	"Holt",
	//	// ringer -- "La",
	//	"Old Bridge",
	//	"Liberty Township",
	//	"Fairfax Station",
	//	"Wharton",
	//	"Piedmont",
	//	"Harvard",
	//	"Firestone",
	//	"Traverse City",
	//	"Fort Wayne",
	//	"Aberdeen",
	//	"Mechanicsburg",
	//	"Indian River",
	//	"Windsor Locks",
	//	"Lumberton",
	//	"Mountainside",
	//	"Freehold",
	//	"Diamond Bar",
	//	"Kentwood",
	//	"Southbury",
	//	"Lebanon",
	//	"Baldwin Park",
	//	// ringer -- "Oracle Redwood Shores",
	//	// UNUNIQUE, lacks s -- "Redwood Shore",
	//	"Chino Hills",
	//	"Sault Sainte Marie",
	//	"Liberty Lake",
	//	"Argyle",
	//	"Juno Beach",
	//	"Mount Prospect",
	//	"Mesquite",
	//	"Naugatuck",
	//	"Canonsburg",
	//	"Ayer",
	//	"Secausus",
	//	"Tinton Fall",
	//	"Hood River",
	//	"Northfield",
	//	"Bethlehem",
	//	"Fresno",
	//	"Ridgefield Park",
	//	"Clara",
	//	"Greencastle",
	//	"Skillman",
	//	"Hackensack",
	//	"Chevy Chase",
	//	"Woodcliffe Lake",
	//	"Surprise",
	//	"Highland Park",
	//	"Sierra Madre",
	//	"Pine River",
	//	"San Carlos City",
	//	"Rosslyn",
	//	"Desert Hot Springs",
	//	"Doylestown",
	//	"Alviso",
	//	"Lyndhurst",
	//	"Eagan",
	//	"Council Bluffs",
	//	"Pelham",
	//	"Mazeppa",
	//	// omit, hides a city -- "Avenel",
	//	"Secaucus",
	//	"Harpers Ferry",
	//	"New Paltz",
	//	"Auburndale",
	//	"Totowa",
	//	"Belle Mead",
	//	"Glenwood",
	//	"Palmyra",
	//	"Evansville",
	//	"Berkeley Heights Township",
	//	"Shrewsbury",
	//	"Valley",
	//	"West Jordan",
	//	// last minute ringer -- "SF",
	//	"Redding",
	//	"Rialto",
	//	"Chesapeake Beach",
	//	"Boaz",
	//	"Loomis",
	//	"Duvall",
	//	"Rosenberg",
	//	// ringer -- "Santa",
	//	"Baldwin",
	//	"Solon",
	//	"Irwindale",
	//	"West Liberty",
	//	"Battery Park",
	//	"Creston",
	//	"Huntington",
	//	"Pepperell",
	//	"Kingwood",
	//	"Newton Highlands",
	//	"Boxborough",
	//	"Fall City",
	//	"Sun Prairie",
	//	"Fairmont",
	//	"Bound Brook",
	//	"Manhasset",
	//	"Voorhees",
	//	"Warrenville",
	//	"Hilton",
	//	"Pooler",
	//	"Weston",
	//	// ringer -- "Global",
	//	"Anchorage",
	//	"Biloxi",
	//	"Shreve",
	//	// UNUNIQUE, lacks s -- "Redwood shore",
	//	"Porter Ranch",
	//	"Jonesboro",
	//	"Peachtree City",
	//	"Post Falls",
	//	"Holden",
	//	// ringer -- "San DIego",
	//	"Haydenville",
	//	"Valley Village",
	//	// UNUNIQUE -- "Jersey city",
	//	"Arnold",
	//	"Rutherford",
	//	"Deltona",
	//	// Rid in 3rd vetting... "Metro West", // Massachusetts Metro West
	//	"Phildelphia",
	//	"Clinton",
	//	"Lakeland",
	//	"Walnut",
	//	"Muskegon",
	//	// ringer -- "DFW",
	//	"Sudbury",
	//	"Fallbrook",
	//	"Nanuet",
	//	"Eureka",
	//	"Carson City",
	//	"Oxnard",
	//	"Eden Praire",
	//	// UNUNIQUE -- "Palo alto",
	//	"Tuckahoe",
	//	"Mukilteo",
	//	"Larkspur",
	//	"Montebello",
	//	"Lorton",
	//	// UNUNIQUE, lacks space -- "Highpoint",
	//	"Erie",
	//	"Collegeville",
	//	"Grand Forks",
	//	"East Elmhurst",
	//	"Belvedere Tiburon",
	//	"Farmingdale",
	//	"New London",
	//	"Willow Grove",
	//	"Rancho Bernardo",
	//	"Phillipsburg",
	//	"Herriman",
	//	// UNUNIQUE -- "Tysons corner",
	//	"Pismo Beach",
	//	"Imperial Beach",
	//	// ringer -- "Silicon Valley Lab San Jose",
	//	"Southborough",
	//	"Branford",
	//	"Santa Fe",
	//	"Northford",
	//	"Ho Ho Kus",
	//	"Union Bridge",
	//	"Vicksburg",
	//	"Oakley",
	//	"Honolulu",
	//	"Tracy",
	//	"Sunnyside",
	//	"Rock Hill",
	//	"Wimberley",
	//	"Roslyn",
	//	// ringer -- "Sanjose",
	//	"San Pablo",
	//	"Pittsburg",
	//	"Pleasonton",
	//	"North Reading",
	//	"Saint Michael",
	//	"Mount Pleasant",
	//	"Berlin",
	//	"Deerfield",
	//	// ringer -- "Austin  Regional Office",
	//	"Glen Gardner",
	//	"Upland",
	//	"Neenah",
	//	"Scottsdate",
	//	"Reynoldsburg",
	//	"Foristell",
	//	"Hillsdale",
	//	"Lake Elsinore",
	//	"Vail",
	//	"Rowland Heights",
	//	"South Bound Brook",
	//	"Bennett",
	//	"Chanhassen",
	//	"Monroe",
	//	"Ellensburg",
	//	"Hilliard",
	//	// ringer -- "Kennedy Space Center",
	//	"Brookline",
	//	"Maplewood",
	//	"Laguna Niguel",
	//	"Cranford",
	//	// ringer -- Lcase, and is special -- "New york",
	//	"Alamo",
	//	"Mims",
	//	"Ft Lauderdale",
	//	"Rogers",
	//	"Croton On Hudson",
	//	"Murray Hill",
	//	"Mcdonough",
	//	"Mount Laurel",
	//	"El Dorado Hills",
	//	"Chatsworth",
	//	"Morris Plains",
	//	"Goleta",
	//	// UNUNIQUE -- "San ramon",
	//	"Fair Lawn",
	//	"Norwood",
	//	"Menomonee Falls",
	//	"Acton",
	//	"Capitol Heights",
	//	"Maricopa",
	//	"Lynbrook",
	//	"Bel Air",
	//	"Northborough",
	//	"Pleasant Prairie",
	//	"Davidsonville",
	//	"Cedar Creek",
	//	"Taylorsville",
	//	"Granada Hills",
	//	"Edinburg",
	//	"Burleson",
	//	"Carefree",
	//	"National City",
	//	"Lilburn",
	//	"Dresher",
	//	"Coalport",
	//	"Portage",
	//	"Uniondale",
	//	"Maywood",
	//	"Franklin Square",
	//	"Lutherville Timonium",
	//	"New Orleans",
	//	"Alexander",
	//	"Desert Oasis",
	//	// ringer -- "Office in Seattle",
	//	"Merchantville",
	//	"Brookfield",
	//	"Braintree",
	//	"Plantation",
	//	// ringer -- "UTC San Diego",
	//	// ringer -- "Santaclara",
	//	"Tappan",
	//	"Plattsburgh",
	//	// 3rd vetting -- "Eastern",
	//	"Goose Creek",
	//	"Marina",
	//	"Orland Park",
	//	"East Greenwich",
	//	"Flagstaff",
	//	"Westcliffe",
	//	"Yorktown",
	//	"Gambrills",
	//	"Hopewell",
	//	"Woburn",
	//	"Storrs",
	//	// ringer -- "REMOTE from Dallas  to Nashville",
	//	"Apple Valley",
	//	"Johnson City",
	//	"New Port Richey",
	//	"North Easton",
	//	"Cocoa",
	//	"North Richland Hills",
	//	"Downingtown",
	//	// ringer -- "Remote in Carlsbad",
	//	"Edmonds",
	//	"Dobbs Ferry",
	//	"Taylor",
	//	"Hyattsville",
	//	"Bedford Hills",
	//	"Hesperia",
	//	// UNUNIQUE, lacks space -- "Menlopark",
	//	"Crown Point",
	//	"Morton Grove",
	//	"Snohomish",
	//	"Clarksburg",
	//	"Mohegan Lake",
	//	"Maryville",
	//	"Lansdale",
	//	"Fortville",
	//	"Saint Joseph",
	//	"Forest Hill",
	//	"Exton",
	//	"Damascus",
	//	"Manasquan",
	//	"Whitestone",
	//	"Lancaster",
	//	"Loveland",
	//	"Catharpin",
	//	"East Rutherford",
	//	"Benson",
	//	"Mead",
	//	"Kearney",
	//	"Smithfield",
	//	"Parlin",
	//	"South Ozone Park",
	//	"Cresskill",
	//	"Menifee",
	//	"North Billerica",
	//	"Bayside",
	//	"Arcadia",
	//	"Dillsburg",
	//	"Prospect Heights",
	//	"North Palm Beach",
	//	"Geneva",
	//	"Mankato",
	//	// ringer -- "Greater Los Angeles",
	//	"Little Silver",
	//	"Plainview",
	//	"New Rochelle",
	//	"Ballston Lake",
	//	"South Pasadena",
	//	// ringer -- "South",
	//	"Jasper",
	//	"Gilcrest",
	//	"Indianola",
	//	"Sound Beach",
	//	"Northport",
	//	"Estero",
	//	"Sun City",
	//	"Bay Shore",
	//	"Woodinville",
	//	"Glen Oaks",
	//	"Bouder",
	//	"Newton Center",
	//	"Elgin",
	//	"Brainerd",
	//	"Matawan",
	//	"Medfield",
	//	"Mountville",
	//	"Annapolis Junction",
	//	"Easton",
	//	"Flower Mound",
	//	"Gig Harbor",
	//	"Apopka",
	//	"Oradell",
	//	"Mchenry",
	//	"Algonquin",
	//	"Merrick",
	//	"Medway",
	//	"Jericho",
	//	"Seaside",
	//	"Monterey Park",
	//	"Franklin Park",
	//	"Fanwood",
	//	"Santa Moncia",
	//	"Swanton",
	//	"Downers Grove",
	//	"Lenexa",
	//	"Ewing",
	//	"Oak Lawn",
	//	"Port Washinton",
	//	"Paramount",
	//	"North Babylon",
	//	"Troy",
	//	"Strongsville",
	//	"Bisbee",
	//	"Chestnut Hill",
	//	"Port Arthur",
	//	"Gibsonville",
	//	"Rapid City",
	//	"Boonsboro",
	//	"Pleasant Valley",
	//	"Mt Laurel",
	//	"Lompoc",
	//	"Camden",
	//	"Scarsdale",
	//	"Bellingham",
	//	"Rowley",
	//	"Eatontown",
	//	"Cedar Park",
	//	// ringer -- "EL Segundo",
	//	"Romeoville",
	//	"Altamonte Springs",
	//	"American Canyon",
	//	"Brookhaven",
	//	"Temple Terrace",
	//	"Starkville",
	//	"Danvers",
	//	"Irvington",
	//	"Rancho Santa Fe",
	//	"Avon",
	//	"Olathe",
	//	"Playa Del Rey",
	//	"Canandaigua",
	//	"Paris",
	//	"Riverton",
	//	"Carson",
	//	"Bay area",
	//	"Greeneville",
	//	"Franktown",
	//	"Holtsville",
	//	"Lake Worth",
	//	"Lees Summit",
	//	"Kensington",
	//	"Los Alamos",
	//	"Northern",
	//	"Winston Salem",
	//	"Irwin",
	//	"Monument",
	//	"Payson",
	//	"Roanoke",
	//	"Huachuca City",
	//	"East Freetown",
	//	"Montrose",
	//	"South Jordan",
	//	"Dunwoody",
	//	"Wayzata",
	//	"Temple City",
	//	"Gardiner",
	//	"Lavergne",
	//	// 3rd vetting -- "Ada", // matches in Canada, okay if rid first; But also matches NEVADA!
	//	"Mission Viejo",
	//	"Sylva", // matches in Pennsylvania; s/b Sylva, North Carolina. Try again with \b
	//	"North Arlington",
	//	// UNUNIQUE, lacks space -- "MountainView",
	//	"Jolla",
	//	"Lake Elmo",
	//	"Lake Stevens",
	//	"Fort Bragg",
	//	// last minute ringer -- "Bay",
	//	"Julian",
	//	"Elmwood Park",
	//	"Leander",
	//	"La Crosse",
	//	"Ogden",
	//	"Houghton",
	//	"Yuba City",
	//	"Denville",
	//	"Roslyn Heights",
	//	"Bellaire",
	//	"Vernon",
	//	"South Grafton",
	//	"Moline",
	//	"West Des Moines",
	//	"Heathrow",
	//	"Centennial",
	//	"Selden",
	//	"Hopewell Junction",
	//	"Mullica Hill",
	//	"Breinigsville",
	//	"Mayfield",
	//	"Little Chute",
	//	"Mobile",
	//	"Eastgate",
	//	"Marion",
	//	// 3rd vetting - ringer -- "Research Lab",
	//	"Elmsford",
	//	"Mount Airy",
	//	"Fort Myers",
	//	"Fuquay Varina",
	//	"Sylmar",
	//	"Castro Valley",
	//	"Vernon Hills",
	//	// UNUNIQUE -- "Westlake village",
	//	"Mid Town",
	//	"Rowlett",
	//	"Malta",
	//	"Stony Point",
	//	"Boynton Beach",
	//	"Phoneix",
	//	"Wesley Chapel",
	//	"Newton",
	//	"Devon",
	//	"Glenside",
	//	"Hayden",
	//	"Stonington",
	//	"Downey",
	//	"Bay Center",
	//	"Waukesha",
	//	"Cheyenne",
	//	"Puyallup",
	//	"Owings Mills",
	//	"Panama City Beach",
	//	"Prescott",
	//	"Elizabeth",
	//	"Quantico",
	//	"West Los Angeles",
	//	"Blue Bell",
	//	"San Deigo",
	//	"Finksburg",
	//	"Edgewater",
	//	"Bloomsburg",
	//	"Freeport",
	//	"Babylon",
	//	"Great Neck",
	//	"Holyoke",
	//	"Bristol",
	//	"Elmira",
	//	"Clarksville",
	//	"Kalamazoo",
	//	"Fountain Hills",
	//	"Bend",
	//	"Takoma Park",
	//	"Clark",
	//	"Paia",
	//	"Stillwater",
	//	"Marblehead",
	//	"Yorba Linda",
	//	"Palos Verdes Peninsula",
	//	"Lima",
	//	"Lithia Springs",
	//	"North Bend",
	//	"Port Orange",
	//	"Colonia",
	//	"Elizabethtown",
	//	"Fair Oaks",
	//	"West Newton",
	//	// UNUNIQUE, lacks space -- "Mountainview",
	//	"Wallingford",
	//	"Vestal",
	//	"Clearwater Tampa",
	//	"Tonawanda",
	//	"Lanham",
	//	"West Hills",
	//	"La Verne",
	//	"Corte Madera",
	//	"Century City",
	//	"Powder Springs",
	//	// ringer -- "San jose",
	//	"West Newbury",
	//	"Park City",
	//	"Auburn Hills",
	//	"Hauppauge",
	//	"Navarre",
	//	"San Gabriel",
	//	"Cortlandt Manor",
	//	"Hacienda Heights",
	//	"Oldsmar",
	//	"Walkersville",
	//	"Lakewood",
	//	"Beacon Falls",
	//	"Waddell",
	//	"Westfield",
	//	"Spring",
	//	// ringer -- "SanFrancisco",
	//	"Jackson Heights",
	//	"Chattanooga",
	//	"Litchfield Park",
	//	"Gerrardstown",
	//	"West Columbia",
	//	"Bristow",
	//	"Wylie",
	//	"Rohnert Park",
	//	"Lakeville",
	//	"Fontana",
	//	"Sykesville",
	//	"Mableton",
	//	"Spokane",
	//	"Bensalem",
	//	"Brentwood",
	//	"Westwood",
	//	"New Bedford",
	//	"Flushing",
	//	"Corvallis",
	//	"Holliston",
	//	"Bartlett",
	//	"Hazlet",
	//	"Floresville",
	//	"Covina",
	//	// ringer -- "DC",
}

type Pair struct {
	key string
	val int
}

type PairList []Pair

func (pl PairList) Len() int      { return len(pl) }
func (pl PairList) Swap(a, b int) { pl[a], pl[b] = pl[b], pl[a] }

// less is more
func (pl PairList) Less(a, b int) bool { return pl[a].val > pl[b].val }

func PrintSortedByFrequency(pq map[string]int, minVal int) {
	pl := PairList{}
	for k, v := range pq {
		if v >= minVal {
			pl = append(pl, Pair{k, v})
		}
	}
	sort.Sort(pl)
	for _, p := range pl {
		fmt.Printf("%5d %s\r\n", p.val, p.key)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var reCRLFs = regexp.MustCompile(`(\r\n|\r|\n)+`)

// create a regexp to strip initial local addresses
// to avoid them matching countries, states, cities
// Reorder -- Else "St." here prevents "St. Cities"
// var reLocalAddrs = regexp.MustCompile(`(?i)^.*\b(avenue|blvd|drive|parkway|street|st\.,|ave\.|university)\b,* *`)
// re-do:
// dropping the puncts after st, ave, oxymoron with \b
// Note, this does lose one city: "University Park, IL", but worth it.
// Also adding dr, rd, s/b okay within \b
var reLocalAddrs = regexp.MustCompile(`(?i)^.*\b(avenue|blvd|drive|dr|rd|parkway|street|st|ave|university)\b,* *`)

// create a regexp to match non-state names
// only in ex_main. new main() does not use
var reNonStates *regexp.Regexp

// That was a gross list. Now just catch the best of DC and NYC:
var reWashDC = regexp.MustCompile(`(?i)(Washington,? ?D\.?C\.?|District Of Columbia|DC Metro)`)
var reNYC = regexp.MustCompile(`(?i)(New York City|New York Area|New York, New York|N(Y|y),? ?NY|New York,? ?NY|NYC)`)

// As with India, amplify \b with OR $: (\b|$)
var reAnyNYorWashState = regexp.MustCompile(`(?i)(New York|\bNY(\b|$)|Washington|\bWA(\b|$))`)

// Make another Ft. and St. ( + specific cities) matcher:
var reFtSt = regexp.MustCompile(`(?i)(Ft\.? (Huachuca|Lauderdale|Worth))|(St\.? (Paul|Louis|Cloud|Joseph|Petersburg))`)

// create a regexp to match state names,
// another regexp for all abbreviations.
// Also make a map from Name to digraph.
var reStateAbbrs *regexp.Regexp
var reStateNames *regexp.Regexp
var stateNameToAbbr = make(map[string]string)

// create a regexp to match country names
var reCountries *regexp.Regexp

// exploratory, only used in ex_main
var goodNonStateQty = map[string]int{}
var goodStateNameQty = map[string]int{}
var goodStateAbbrQty = map[string]int{}
var goodCountryQty = map[string]int{}
var badResidueQty = map[string]int{}
var qualityCityQty = map[string]int{}

// Golden. in main2
var reQualityCities *regexp.Regexp

// More Golden. in main3
var rePerfected *regexp.Regexp

func initTools() {
	// create a regexp to match non-state names
	{
		sb1 := []byte{}
		for _, ns := range nonStates {
			sb1 = append(sb1, byte('|'))
			sb1 = append(sb1, []byte(ns)...)
		}
		sb1 = append(sb1, byte(')'))
		sb1[0] = '(' // overwrite first pipe
		reNonStates = regexp.MustCompile(string(sb1))
	}

	// create a regexp to match state names,
	// another regexp for all abbreviations.
	// Also make a map from Name to digraph.
	{
		sb0 := []byte{byte('\\'), byte('b')}
		// late in the game, but add (?i) prefix on Name, not Abbr:
		// Nice, but rather: Need the india trick here ATOP line:
		// No, I do not have a top of line circumflex here...
		sb1 := []byte{'(', '?', 'i', ')'}
		for _, AbbrName := range states {
			an := strings.Split(AbbrName, ",")
			stateAbbr := an[0]
			stateName := an[1]
			sb0 = append(sb0, byte('|'))
			sb0 = append(sb0, []byte(stateAbbr)...)
			sb1 = append(sb1, byte('|'))
			sb1 = append(sb1, []byte(stateName)...)
			stateNameToAbbr[stateName] = stateAbbr
		}
		// Need the india trick here too, both places:
		// Also insert the close paren to alternation:
		sb0 = append(sb0, []byte(`)(\b|$)`)...)
		sb1 = append(sb1, []byte(`)(\b|$)`)...)
		sb0[2] = '(' // overwrite first pipe
		sb1[4] = '(' // overwrite first pipe
		reStateAbbrs = regexp.MustCompile(string(sb0))
		reStateNames = regexp.MustCompile(string(sb1))
	}

	// create a regexp to match country names
	{
		sb1 := []byte{}
		for _, cn := range countries {
			sb1 = append(sb1, byte('|'))
			sb1 = append(sb1, []byte(cn)...)
		}
		sb1 = append(sb1, byte(')'))
		sb1[0] = '(' // overwrite first pipe
		reCountries = regexp.MustCompile(string(sb1))
	}

	// more in new main...

}

func ex_main() {

	ba, err := ioutil.ReadAll(os.Stdin)
	check(err)
	fmt.Println("Bytes", len(ba))

	sa := reCRLFs.Split(string(ba), -1)
	fmt.Println("Lines", len(sa))

	initTools()

	// Finally, loop over input lines

	for _, s := range sa {

		s = reLocalAddrs.ReplaceAllString(s, "")

		if reNonStates.Match([]byte(s)) {
			goodNonStateQty[s] = goodNonStateQty[s] + 1
			continue
		}

		if reStateNames.Match([]byte(s)) {
			goodStateNameQty[s] = goodStateNameQty[s] + 1
			developCities(reStateNames.ReplaceAllString(s, ""), &qualityCityQty)
			continue
		}

		if reStateAbbrs.Match([]byte(s)) {
			goodStateAbbrQty[s] = goodStateAbbrQty[s] + 1
			developCities(reStateAbbrs.ReplaceAllString(s, ""), &qualityCityQty)
			continue
		}

		if reCountries.Match([]byte(s)) {
			goodCountryQty[s] = goodCountryQty[s] + 1
			continue
		}

		badResidueQty[s] = badResidueQty[s] + 1

	}

	/*
		fmt.Println("goodNonStateQty", len(goodNonStateQty))
		PrintSortedByFrequency(goodNonStateQty, 5)

		fmt.Println("goodStateNameQty", len(goodStateNameQty))
		PrintSortedByFrequency(goodStateNameQty, 5)

		fmt.Println("goodStateAbbrQty", len(goodStateAbbrQty))
		PrintSortedByFrequency(goodStateAbbrQty, 5)

		fmt.Println("goodCountryQty", len(goodCountryQty))
		PrintSortedByFrequency(goodCountryQty, 5)

	*/

	fmt.Println("badResidueQty", len(badResidueQty))
	PrintSortedByFrequency(badResidueQty, 5)

	// After many reductions, this list of quality cities
	// will become a most important input clustering tool

	fmt.Println("qualityCityQty", len(qualityCityQty))
	PrintSortedByFrequency(qualityCityQty, 1)

	for _, s := range rejected {
		fmt.Println("REJECTED", s)
	}
}

// this is all ex_main stuff...

var reSlashSplitter = regexp.MustCompile(` */ *`)

// this has grown some!
var reSpCommaSpEnd = regexp.MustCompile(`(^Greater )|((-| |,|\.|Area|USA|US|\d+|Estados Unidos|United States|U.S.A)+$)`)

var reUndesirables = regexp.MustCompile(`[[:^ascii:][:digit:][:punct:]]`) // -[-]

var rejected = make([]string, 0)

// this is all ex_main stuff...

func developCities(s1 string, qualityCityQty *map[string]int) {

	// caller stripped the StateName or StateAbbr

	// Now strip any final comma, spaces
	// Also Greater ... Area parts, and USA
	// then tabulate remaining hi-Q CityNames

	// Loop over '/' separating several names

	// Fix initial Ft. or St. to Fort and Saint.
	// Reject items with any other digits,
	// puncts save '-', non-UC[0], non-usascii

	sa := reSlashSplitter.Split(s1, -1)
	for _, s := range sa {
		s = reSpCommaSpEnd.ReplaceAllString(s, "")
		s = strings.Replace(s, "Ft.", "Fort", 1)
		s = strings.Replace(s, "St.", "Saint", 1)
		// with all that cleaned up, split on '-'
		sa2 := strings.Split(s, "-")
		for _, s2 := range sa2 {
			s2 = strings.TrimSpace(s2)
			if len(s2) > 0 &&
				unicode.IsUpper([]rune(s2)[0]) &&
				reUndesirables.Match([]byte(s2)) == false {
				(*qualityCityQty)[s2] = (*qualityCityQty)[s2] + 1
			} else {
				rejected = append(rejected, fmt.Sprintf("REJECTED [[%s]] in [[%s]]", s2, s1))
			}
		}
	}
}

// QC = Quality City names

// turns out there were no QC in >1 state, unneeded:
// var lcQC2ListStateAbbr = make(map[string][]string)
// I was wrong about that, but recovered, guess a few

var lcQC2StateAbbr = make(map[string]string)
var lcQCBareNoSA = make(map[string]string)

/* Great! Just 5 Bare QC; Insert these into ... somewhere...:
Chesterbrook, PA
Lilburn, GA
Trenton, NJ
West Lafayette, IN
Hopewell, VA
*/

func ex_main_two() {
	// I already passed Go, collected $200.
	// Now I can use the QualityCities list
	// that was generated by ex_main above.

	initTools()

	// Study QualityCities just within itself.
	// done...
	// {
	// 	unique := make(map[string]bool)
	// 	for _, s := range QualityCities {
	// 		s = strings.ToLower(s)
	// 		// repeat without spaces
	// 		s = strings.Replace(s, " ", "", -1)
	// 		if _, ok := unique[s]; ok {
	// 			fmt.Println("UNUNIQUE", s)
	// 		}
	// 		unique[s] = true
	// 	}
	// }

	// Prepare a regexp, longest strings first:
	{
		sort.Sort(QualityCities)
		// throw in a case insensitive (?i) prefix:
		sb1 := []byte{'(', '?', 'i', ')'}
		for _, cn := range QualityCities {
			sb1 = append(sb1, byte('|'))
			sb1 = append(sb1, []byte(cn)...)
		}
		sb1 = append(sb1, byte(')'))

		// Optional:
		// adding \b to omit Brea From Panera Bread
		// and to fix 2-3 others restored to QC list
		sb1 = append(sb1, byte('\\'))
		sb1 = append(sb1, byte('b'))

		sb1[4] = '(' // overwrite first pipe
		// prove the sort:
		// good... fmt.Println(string(sb1))
		reQualityCities = regexp.MustCompile(string(sb1))
	}

	// Study City + State over input data
	ba, err := ioutil.ReadAll(os.Stdin)
	check(err)
	sa := reCRLFs.Split(string(ba), -1)
	for _, line := range sa {

		bline := []byte(line)

		// TO DO: each one-off that continues
		// must invoke eventual line output provision...

		// This is a one-off, City implies State:
		// Must do before ridding St. etc prefix.
		if reFtSt.Match(bline) {
			// fmt.Println("Fort/Saint here:", line)
			continue
		}

		// Run this before QualityCities match,
		// prevents many false MULTIPLE city matches.
		line = reLocalAddrs.ReplaceAllString(line, "")

		bline = []byte(line) // again

		// This is essentially all foreign gigo:
		// Run before Quality Cities to fix a few.
		if reCountries.Match(bline) {
			continue
		}

		// This is a one-off, City implies State:
		if reNYC.Match(bline) {
			// fmt.Println("NYC here:", line)
			continue
		}

		// This is a one-off, City implies State:
		// Run before Quality Cities' "Columbia".
		if reWashDC.Match(bline) {
			// fmt.Println("WashDC here:", line)
			continue
		}

		if reQualityCities.Match(bline) {
			// wip... fmt.Println(line)

			// I want to study any State Name or Abbr right of City:
			// in fact, all residue either side, and multiplicities.
			ss := reQualityCities.Split(line, -1)
			switch len(ss) {
			// explanation: Match before Split pbly prevented 0, 1:
			case 0:
				// never seen
				// perhaps City was alone on line?
				// fmt.Println("ZERO:", ss)
				break
			case 1:
				// never seen
				// perhaps City was atop line
				// fmt.Println("ONE: ", ss)
				break
			case 2:
				// Normally expected condition, exactly one City:
				// I also fall into 2 for no text left of a City.
				// I think that insight is how golang Split works!
				// fmt.Println(ss)

				// So, in here, study STATE NAME|ABBR in ss[1]:

				// But first, dig out the QualityCity match:
				qc := string(reQualityCities.Find(bline))
				// fmt.Println("QC:", qc)

				// Use LC for indexing
				lcqc := strings.ToLower(qc)

				bss1 := []byte(ss[1])
				stateAbbr := ""
				if reStateNames.Match(bss1) {
					ssn := string(reStateNames.Find(bss1))
					stateAbbr = stateNameToAbbr[ssn]
					// great... fmt.Println(stateAbbr)
				} else {
					if reStateAbbrs.Match(bss1) {
						stateAbbr = string(reStateAbbrs.Find(bss1))
						// great... fmt.Println(stateAbbr)
					}
				}

				if stateAbbr == "" {
					// have NO State for this City
					// fmt.Println("QC ALONE:", qc)

					// No state abbr:
					lcQCBareNoSA[lcqc] = ""

				} else {
					// have a State for this City
					// qcsa := qc + ", " + stateAbbr
					// fmt.Println(qcsa)
					// Tabulate to see if any qc cross over sa

					// none, change to a single sa:
					// lsa, ok := lcQC2ListStateAbbr[lcqc]
					// if ok {
					// 	lsa = append(lsa, stateAbbr)
					// } else {
					// 	lcQC2ListStateAbbr[lcqc] = []string{stateAbbr}
					// }

					// single abbr:
					lcQC2StateAbbr[lcqc] = stateAbbr

				}

				break
			default:
				// perhaps multiple cities?
				// fmt.Println("MULTIPLE: ", ss)

				// N.B. This main() has NOT
				// split on '/' as did ex_main.

				// TO DO: Decide what to do.
				// Take first, last, ignore?
				break
			}

		} else {
			// reduce the residue (non-QC) to study:

			// fmt.Println(line)
		}
	}

	// after all input lines:

	// study lcQC2ListStateAbbr for 1, 2+
	// Good! Done. There were no cities with wrong states!
	// for qc, lsa := range lcQC2ListStateAbbr {
	// if len(lsa) > 1 {
	// fmt.Println(qc, lsa)
	// }
	// }
	// Oops, my logic failed me somewhere there.
	// After ridding the SA list code provision,
	// 2 runs got different results on a few names,
	// (thank goodness Golang randomizes map order)
	// So I manually added them to perfected list.

	// What QC remain with no SA?
	// for qc, _ := range lcQCBareNoSA {
	// 	if _, ok := lcQC2StateAbbr[qc]; !ok {
	// 		fmt.Println("NoSa:", qc)
	// 	}
	// }
	// Great! Just 5 to google; Insert these somewhere:
	// Chesterbrook, PA
	// Lilburn, GA
	// Trenton, NJ -- this one came in next time.
	// West Lafayette, IN
	// Hopewell, VA

	// Oh, I don't have any somewhere...
	// Create a new QC list with SA...
	// Done. This + sorting + a few lines => perfected list.
	// for qc, sa := range lcQC2StateAbbr {
	// 	qc = strings.ToLower(qc)
	// 	qc = strings.Title(qc) // despite deprecated; not .ToTitle
	// 	qc = qc + ", " + sa
	// 	fmt.Println(qc)
	// }
}

// ============= go on to make a third main: ===========

var bareTcCity2Perfect = map[string]string{}
var bareTcCityDuplicate = map[string]bool{}

// Plug a few last city sieve holes here:
// These all report clue, (No City Match) + 2 crazy
// NYC ...
// DC ...
// DFW ...
// John's Creek, GA
// Tyson's Corner, VA
// Ada, Oklahoma
// University Park, Il -- must go back to virgin line
// But omit just a few actual universities that dropped out
// == NYC|DC|DFW|John's Creek|Tyson's Corner|Ada,|University Park|

// and perhaps allow:
// RTP, NC => Raleigh, NC
// Research Triangle Park, NC
// NewYork
// Norwell,Ma
// Hawaiian Islands
// IBM Almaden... = SanJose CA
// IBM Thomas... = Yorktown Heights NY
// SFO
// SF
// But not Portsmouth, too many states
// == RTP|Research Triangle|NewYork|Norwell|Hawaiian|IBM |SF|

//We love engineers:
//NASA Ames Research Center => Mountain View, CA
//Silicon Valley Lab... => San Jose, CA
//Facebook => Menlo Park, CA
//Kennedy Space Center... => Merritt Island, ‎FL
// == NASA Ames|Silicon Valley Lab|Facebook|Kennedy Space Center|

//Engin'rs cant spell:
//Sanjose
//Redwood shore
//Mountainview
//Menlopark
//SanFrancisco
//Santa Moncia, CA
// == Sanjose|Redwood shore|Mountainview|Menlopark|SanFrancisco|Santa Moncia

// and more latecomers: These should put the cherry on top!

// But not all, keep only if atop line:

// NASA Headquarters -> dc
// Cambidge Ma
// El Segund
// Glendora -- add to city list
// Greeley -> co ditto
// McAfee --> Plano, TX -- but use Intel Security / McAfee
// Jersey - skip, too many places
// Malibu - I put that in, what happened? Oh, foreign? - Oh, Mali!
// National Cybersecurity FFRDC -> Rockville, MD
// Portsmouth - skip, too many places
// Robins AFB - Robins, GA
// SAIC -- skip, too many
// SFO - still MIA, ahh, not matched by SF\b
// San Diegp -- these should just come in to switch after adding to the regexp
// San Ramon, CA -- add to list. It is there. Hmmm.
// SanFransisco
// Sen Diego
// Virigina
// WFH? - Ah, WorkFromHome! And WW = WorldWide.
// ca, pleasanon
// california -- a mystery
// rtp -- now lowercase
// thousand oak
//

// NASA Headquarters|Cambidge Ma|El Segund|Intel Security / McAfee|National Cybersecurity FFRDC|Robins AFB|San Diegp|SanFransisco|Sen Diego|Virigina|ca, pleasanon|thousand oak|

// I would not think a comma goes under the \b span, but see if it fixes "MountainView,CA". No. Capital V!
// But do move the , from Ada, to end.
// the reason I need the well-spelled San Francisco, is it was saved from foreigns
var rePlugSieve = regexp.MustCompile(`^(NYC|DC|DFW|John's Creek|Tyson's Corner|Ada|University Park|RTP|rtp|Research Triangle|NewYork|Norwell|Hawaiian|IBM |SF|SFO|NASA Ames|Silicon Valley Lab|Facebook|Kennedy Space Center|(S|s)an(J|j)ose|Redwood (S|s)hore|Mountain(v|V)iew|Menlopark|SanFrancisco|Santa Moncia|NASA Headquarters|Cambidge Ma|El Segund|Intel Security / McAfee|National Cybersecurity FFRDC|Robins AFB|San Diegp|SanFransisco|Sen Diego|Virigina|ca, pleasanon|thousand oak|San Francisco)(\b|,|$)`)

func main() {
	// I passed Go twice, collected $400.
	// Now I can use the "perfected" list
	// that was generated by ex_main_two.

	initTools()

	// Prepare a regexp, longest strings first:
	// But this time QualityCities = "perfected" list.
	// NOW, I must strip off a comma space state abbr.
	{
		// N.B. this is a sort by descending string length:
		sort.Sort(perfected)
		// throw in a case insensitive (?i) prefix:
		sb1 := []byte{'(', '?', 'i', ')'}

		// Now working in Title Case cities + ", XX"
		for _, tcCity2Abbr := range perfected {

			sb1 = append(sb1, byte('|'))
			cityOnly := []byte(tcCity2Abbr)
			cityOnly = cityOnly[:len(cityOnly)-4]
			sb1 = append(sb1, cityOnly...)

			// Just once, verify I still have all TC:
			// good...
			// chk := string(cityOnly)
			// chk2 := strings.Title(strings.ToLower(chk))
			// if chk2 != chk {
			// 	fmt.Println("PLEASE REPAIR:", chk)
			// }

			// Meanwhile, build up fresh translation maps
			// 1. for cities with explicitly stated state
			// Actually, in title case, idempotent, not needed.
			// 2. note a few cities in multiple states to study
			// 3. for bare city lookup, choosing just one state
			strTcCityOnly := string(cityOnly)
			if _, ok := bareTcCity2Perfect[strTcCityOnly]; ok {
				bareTcCityDuplicate[strTcCityOnly] = true
			} else {
				// There are a few duplicate cities, overwrite:
				bareTcCity2Perfect[strTcCityOnly] = tcCity2Abbr
			}
		}
		// still finishing that regexp:
		// Without \b I lost all Pennsylvania and Kentucky, about 50 items
		// With \b I saw no detriment. Keeping...
		// No, Rather, as with India, amplify \b with OR $: (\b|$)
		// Also insert the close paren to alternation:
		sb1 = append(sb1, []byte(`)(\b|$)`)...)

		sb1[4] = '(' // overwrite first pipe after (?i)
		rePerfected = regexp.MustCompile(string(sb1))
	}

	// prove me this, once:
	// fmt.Printf("%#v\r\n", bareTcCityDuplicate)
	// good. Finally...
	// map[string]bool{"Bloomfield":true, "Bloomington":true, "Covington":true, "Danville":true, "Dayton":true, "Harrison":true, "Hopewell":true, "Jackson":true, "Louisville":true, "Medford":true, "Montclair":true, "Morristown":true, "Newark":true, "Oakland":true, "Reading":true, "Richmond":true, "Somerville":true, "Union":true, "Union City":true, "Wayne":true, "Woodbridge":true}

	// Third and Final production run over input data
	ba, err := ioutil.ReadAll(os.Stdin)
	check(err)
	sa := reCRLFs.Split(string(ba), -1)
	for _, line := range sa {

		// keep virgin input line for a
		// final analysis and table output
		virgin := line
		bline := []byte(line)
		clues := "" // to diagnose paths taken

		if reFtSt.Match(bline) {
			// this time, revise the line
			// All ft. were "Ft. " or "Ft "
			// All st. were "St. " exactly.
			// All st for street lacked a .
			line = strings.Replace(line, "Ft. ", "Fort ", -1)
			line = strings.Replace(line, "Ft ", "Fort ", -1)
			line = strings.Replace(line, "St. ", "Saint ", -1)
			// not need, done again below...
			// bline = []byte(line) // again
			clues = clues + "(Ft.St.)"
		}

		// Run this before "rePerfected" match,
		// prevents many false MULTIPLE city matches.
		line = reLocalAddrs.ReplaceAllString(line, "")
		bline = []byte(line) // again

		// Rid foreign countries ASAP to fix a few.
		if reCountries.Match(bline) {
			// Do no more continues in main #3
			line = ""            // spoil the line
			bline = []byte(line) // again
			clues = clues + "(Foreign)"
		}
		// No more here... if reNYC.Match(bline) ...
		// No more here... if reWashDC.Match(bline) ...

		answer := ""

		// this is now using the "perfected", not QC list.
		if rePerfected.Match(bline) {

			// A single IF logic got way too convoluted

			doit := false
			gotNyWa := false

			ss := rePerfected.Split(line, -1)
			if len(ss) == 2 {
				doit = true
				clues = clues + "(==2)"
			} else if len(ss) > 2 {
				clues = clues + "(>2)"

				// Above original split -1 cut out
				// the 2nd city/state text. Re-do:
				ss2 := rePerfected.Split(line, 2)
				if reAnyNYorWashState.Match([]byte(ss2[1])) {
					doit = true
					gotNyWa = true
					clues = clues + "(NyWa)"
				} else {
					// rehabilitate twin cities,
					// and multiple city lists,
					// by just doing it...

					// That introduced 2 crazies,
					// Ames, IA (>2)(Just Do It!) **FOR** NASA Ames Research Center, Moffett Field, Mountain View, CA
					// Silicon Valley, CA (>2)(Just Do It!) **FOR** Silicon Valley Lab San Jose CA 95141
					// which I anally now exclude:
					// Or better, solve:
					tc := string(rePerfected.Find(bline))
					tc = strings.ToLower(tc)
					tc = strings.Title(tc) // N.B. not .ToTitle()!
					if tc == "Ames" {
						clues = clues + "(Crazy)"
						answer = "Mountain View, CA"
					} else if tc == "Silicon Valley" {
						clues = clues + "(Crazy)"
						answer = "San Jose, CA"
					} else {
						doit = true
						clues = clues + "(Just do it!)"
					}
				}
			}

			if doit {
				// Use Title Case of city matched in line for indexing
				tc := string(rePerfected.Find(bline))
				// Also rid ' from John's and Tyson's
				// Didn't help, recover them later.
				// tc = strings.Replace(tc, "'", "", -1)
				tc = strings.ToLower(tc)
				tc = strings.Title(tc) // N.B. not .ToTitle()!

				if !gotNyWa {
					// Solve any other State Name or Abbr right of City:
					bss1 := []byte(ss[1])
					stateAbbr := ""
					if reStateNames.Match(bss1) {
						ssn := string(reStateNames.Find(bss1))
						stateAbbr = stateNameToAbbr[ssn]
					} else {
						if reStateAbbrs.Match(bss1) {
							stateAbbr = string(reStateAbbrs.Find(bss1))
						}
					}
					if stateAbbr == "" {
						// have NO State for this City
						answer = bareTcCity2Perfect[tc]
						if bareTcCityDuplicate[tc] {
							clues = clues + "(Guess)"
						}
					} else {
						// have a State for this City
						answer = tc + ", " + stateAbbr
					}
				} else {
					// I know state for a City in NY or WA
					answer = bareTcCity2Perfect[tc]
					// note it, just in case:
					if bareTcCityDuplicate[tc] {
						clues = clues + "(Guess)"
					}
				}
			} else {
				clues = clues + "(Not doit)"
			}
		} else {
			clues = clues + "(No City Match)"

			solved := false

			// Plug a few last city sieve holes here:
			// SORTED FOR CODING:
			// Ada, Oklahoma
			// DC ...
			// DFW ...
			// Facebook => Menlo Park, CA
			// Hawaiian Islands
			// IBM Almaden... = SanJose CA
			// IBM Thomas... = Yorktown Heights NY
			// John's Creek, GA
			// Kennedy Space Center... => Merritt Island, ‎FL
			// Menlopark
			// Mountainview
			// NASA Ames Research Center => Mountain View, CA
			// NYC ...
			// NewYork
			// Norwell,Ma
			// RTP, NC => Raleigh, NC
			// Redwood shore
			// Research Triangle Park, NC
			// SF
			// SFO
			// SanFrancisco
			// Sanjose
			// Santa Moncia, CA
			// Silicon Valley Lab... => San Jose, CA
			// Tyson's Corner, VA
			// University Park, Il -- must go back to virgin line

			// more added since that list...
			/*
				NASA Headquarters|
				Cambidge Ma|
				El Segund|
				McAfee|
				National Cybersecurity FFRDC|
				Robins AFB|
				San Diegp|
				SanFransisco|
				Sen Diego|
				Virigina|
				ca, pleasanon|
				thousand oak|
			*/
			bvirgin := []byte(virgin)
			if rePlugSieve.Match(bvirgin) {
				//matrix
				switch rune(bvirgin[0]) {
				case 'A':
					answer = "Ada, OK"
					solved = true
					break
				case 'C':
					answer = "Cambridge MA"
					solved = true
					break
				case 'c':
					answer = "Pleasanton, CA"
					solved = true
					break
				case 'D':
					switch rune(bvirgin[1]) {
					case 'C':
						answer = "Washington, DC"
						solved = true
						break
					case 'F':
						answer = "Dallas, TX"
						solved = true
						break
					}
					break
				case 'E':
					answer = "El Segundo, CA"
					solved = true
					break
				case 'F':
					answer = "Menlo Park, CA"
					solved = true
					break
				case 'H':
					answer = "Unknown City, HI"
					solved = true
					break
				case 'I':
					switch rune(bvirgin[4]) {
					case 'A':
						answer = "San Jose, CA"
						solved = true
						break
					case 'l':
						answer = "Plano, TX"
						solved = true
						break
					case 'T':
						answer = "Yorktown Heights, NY"
						solved = true
						break
					}
					break
				case 'J':
					// Drop the apostraphe, for narrow uniformity
					answer = "Johns Creek, GA"
					solved = true
					break
				case 'K':
					answer = "Merritt Island, ‎FL"
					solved = true
					break
				case 'M':
					switch rune(bvirgin[1]) {
					case 'e':
						answer = "Menlo Park, CA"
						solved = true
						break
					case 'o':
						answer = "Mountain View, CA"
						solved = true
						break
					}
					break
				case 'N':
					switch rune(bvirgin[1]) {
					case 'A':
						switch rune(bvirgin[5]) {
						case 'A':
							answer = "Mountain View, CA"
							solved = true
							break
						case 'H':
							answer = "Washington, DC"
							solved = true
							break
						}
					case 'a':
						answer = "Rockville, MD"
						solved = true
						break
					case 'e':
						fallthrough
					case 'Y':
						answer = "New York, NY"
						solved = true
						break
						// Not needed, just added to cities list:
						// case 'o':
						// 	answer = "Norwell, MA"
						// 	solved = true
						// 	break
					}
					break
				case 'r':
					fallthrough
				case 'R':
					switch rune(bvirgin[2]) {
					case 'p':
						fallthrough
					case 'P':
						fallthrough
					case 's':
						answer = "Raleigh, NC"
						solved = true
						break
					case 'd':
						answer = "Redwood Shores, CA"
						solved = true
						break
					case 'b':
						answer = "Robins, GA"
						solved = true
						break
					}
					break
				case 's':
					fallthrough
				case 'S':
					switch rune(bvirgin[1]) {
					case 'i':
						answer = "San Jose, CA"
						solved = true
						break
					case 'F':
						answer = "San Francisco, CA"
						solved = true
						break
					case 'e':
						answer = "San Diego, CA"
						solved = true
						break
					case 'a':
						switch rune(bvirgin[4]) {
						case 'a':
							answer = "Santa Monica, CA"
							solved = true
							break
						case 'o':
							fallthrough
						case 'j':
							answer = "San Jose, CA"
							solved = true
							break
						case 'r':
							fallthrough
						case 'F':
							answer = "San Francisco, CA"
							solved = true
							break
						case 'D':
							answer = "San Diego, CA"
							solved = true
							break
						}
						break
					}
					break
				case 't':
					fallthrough
				case 'T':
					switch rune(bvirgin[1]) {
					case 'h':
						answer = "Thousand Oaks, CA"
						solved = true
						break
					case 'y':
						// Drop the apostraphe, for narrow uniformity
						answer = "Tysons Corner, VA"
						solved = true
						break
					}
					break
				case 'U':
					answer = "University Park, IL"
					solved = true
					break
				case 'V':
					answer = "Unknown City, VA"
					solved = true
					break
				}
				if solved {
					clues = clues + "(solved)"
				} else {
					clues = clues + "(Plug Sieve)"
				}
			} else {
				clues = clues + "(Check Again)"
			}

			if !solved {
				// Better to have just a state name than nothing.

				// Again, Solve any State Name or Abbr:
				// but this time using the whole bline:
				stateAbbr := ""
				// hmmm - not matching lowercase california?
				// perhaps go back to virgin line?
				// Aha, because it was atop/alone on line
				// No, there is not a circumflex in regexp
				// Some mystery I will not waste more time to solve...
				bvirgin := []byte(virgin)
				if reStateNames.Match(bvirgin) {
					ssn := string(reStateNames.Find(bvirgin))
					stateAbbr = stateNameToAbbr[ssn]
				} else {
					if reStateAbbrs.Match(bvirgin) {
						stateAbbr = string(reStateAbbrs.Find(bvirgin))
					}
				}
				if stateAbbr == "" {
					// have NO State for this City
					// but, no city either.
					// Nothing I can do.
				} else {
					// have a State, but no city
					// Name the city Unknown City,
					// and keep the state as Abbr.
					answer = "Unknown City, " + stateAbbr
				}
			}
		}
		if answer == "" {
			// It has been about a week of coding,
			// there will always be more to polish
			// fmt.Println("000 NO ANSWER", clues, "FOR", virgin)

		} else {
			// fmt.Println(answer, clues, "**FOR**", virgin)

			// It is time for production. Make C# code strings,
			// two per line: the translation, then the source.
			// My original DB prevented backslash, dbl quotes.
			fmt.Printf("\t\t\"%v\", \"%v\",\r\n", answer, virgin)
		}
	}
}
