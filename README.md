# Restaurant Search

## Database table creation
```
CREATE TABLE `restaurant` (
  `id` int NOT NULL AUTO_INCREMENT,
  `restaurant_id` int(11) NOT NULL,
  `name` varchar(128) DEFAULT NULL COMMENT 'restaurant name',
  `url` varchar(1024) DEFAULT NULL COMMENT 'restaurant web URL',
  `cuisines` varchar(1024) DEFAULT NULL COMMENT 'comma separated cuisines',
  `image` varchar(1024) DEFAULT NULL COMMENT 'image URL',
  `address` varchar(1024) DEFAULT NULL COMMENT 'complete address',
  `city` varchar(128) DEFAULT NULL COMMENT 'geographics',
  `rating` float DEFAULT NULL COMMENT 'rating max 5',
  `veg` tinyint(1) DEFAULT '0' COMMENT 'type veg or non-veg',
  `createdOn` datetime DEFAULT CURRENT_TIMESTAMP,
  `updatedOn` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
```

## Configuration for Backend

- modify the DB configurations as per the environemnt

```bash
cd restaurant-search-backend
vi conf.yml
```

## Running Backend Server

```bash
  cd restaurant-search-backend

  # Start backend server
  go run restaurant-search-backend

```

The App will run at `localhost:5000`

supported URL & Method:  `localhost:5000/search?query=value1&veg=value2` - GET

**Sample Response**: 

```json
{
  "data": {
    "count": 2,
    "restaurants": [
      {
        "_id": "5f7ebf4b7695e7905ad619ba",
        "restaurant_id": 18353121,
        "name": "Flechazo",
        "url": "https://www.zomato.com/bangalore/flechazo-marathahalli?utm_source=api_basic_user&utm_medium=api&utm_campaign=v2.1",
        "cuisines": "Asian, Mediterranean, North Indian",
        "image": "https://b.zmtcdn.com/data/pictures/1/18353121/884830429217e19190d945fbf9b5351e_featured_v2.jpg",
        "address": "9/1, 1st Floor, Above Surya Nissan, VRR Orchid, Doddanakkundi, Marathahalli, Bangalore",
        "city": "Bangalore",
        "rating": 4.4,
        "veg": true
      },
      {
        "_id": "5f7ebf4b7695e7905ad619bc",
        "restaurant_id": 51040,
        "name": "Truffles",
        "url": "https://www.zomato.com/bangalore/truffles-koramangala-5th-block?utm_source=api_basic_user&utm_medium=api&utm_campaign=v2.1",
        "cuisines": "American, Burger, Cafe",
        "image": "https://b.zmtcdn.com/data/pictures/chains/8/51038/c2c164cf25a35e98f79576691f3a5622_featured_v2.png",
        "address": "28, 4th 'B' Cross, Koramangala 5th Block, Bangalore",
        "city": "Bangalore",
        "rating": 4.7,
        "veg": true
      }
    ]
  }
}
```
