# GPX files location finder

Find gpx files according to the location where they were recorded.

Target location can be :
 - [x] a bounding box
 - [x] a circle 


```
NAME:                                                                                                                     
   gpxloc - GPX files Location Finder                                                                                     
                                                                                                                          
   Find gpx files according to the location where they were recorded                                                      
                                                                                                                          
USAGE:                                                                                                                    
   gpxloc [global options] [arguments...]                                                                                 
                                                                                                                          
VERSION:                                                                                                                  
   0.1                                                                                                                    
                                                                                                                          
GLOBAL OPTIONS:                                                                                                           
   --bbox value                    bounding box ("lon1,lat1,lon2,lat2")                                                   
   --latitude value, --lat value   latitude (WGS84 [-90,+90]) (default: 0)                                                
   --longitude value, --lon value  longitude (WGS84 [-180,+180]) (default: 0)                                             
   --radius value                  radius (in meters) (default: 0)                                                        
   --help, -h                      show help                                                                              
   --version, -v                   print the version                                                                      
                                                                                                                          
                                                                                                                          
WEBSITE: https://github.com/JVillafruela/GPXloc                                                                           
                                                                                                                          
EXAMPLES:                                                                                                                 
   gpxloc --bbox="5.68678,45.08596,5.68979,45.08778" E:\OSM\gps\2022\2022-04-16 E:\OSM\gps\2022\2022-04-22                
                                                                                                                          
   gpxloc --latitude=45.087 --longitude=5.688 --radius=20  E:\OSM\gps\2023                                               
                                                                                      
```