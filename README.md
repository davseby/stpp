Modulio kodas: T120B165
Dėstytojai:
- Doc. Prakt. Rasa Mažutienė
- Doc. Prakt. Petras Tamošiūnas

Studentas:
- Dovydas Bykovas, IFF-9/2 grupė

# Foodie (Saityno taikomųjų programų projektavimas)

Projekto tikslas - viešą maisto receptų bei jų planų kūrimo ir 
dalijimosi platformą.

Veikimo principas - platformą sudaro internetinė sąsaja (frontend) ir serverio 
aplikacija (backend). Internetinės sąsajos ir serverio aplikacija komunikacija 
vyksta per RESTful API.

Sistema turi tris vartotojų grupes:

- Svečias
- Vartotojas
- Administratorius

Svečio funkcijos:

- Užsiregistruoti.
- Peržiūrėti produktus.
- Peržiūrėti pasirinktą produktą.
- Peržiūrėti receptus.
- Peržiūrėti pasirinktą receptą.
- Peržiūrėti vartotojo receptus.
- Peržiūrėti planus.
- Peržiūrėti pasirinktą planą.
- Peržiūrėti vartotojo planus.
- Pasiimti API versiją.

Vartotojo funkcijos:

- Visos svečio funkcijos.
- Prisijungti.
- Sukurti receptą.
- Atnaujinti receptą.
- Ištrinti receptą.
- Sukurti planą.
- Atnaujinti planą.
- Ištrinti planą.
- Ištrinti savo vartotoją.
- Pasikeisti slaptažodį.

Administratorius funkcijos:

- Visos vartotojo funkcijos.
- Sukurti produktą.
- Atnaujinti produktą.
- Ištrinti produktą.
- Ištrinti vartotojo receptą.
- Ištrinti vartotojo planą.
- Ištrinti vartotoją ar administratorių.
- Pasiimti vartotojų sąrašą.
- Sukurti administracinį vartotoją.

# Sistemos architektūra

Aplikacijos frontend dalį sudaro JavaScript ir Vue 3 karkasas.
Aplikacijos backend dalį sudaro Go programavimo kalba, naudojami Docker bei 
docker-compose įrankiai bei MariaDB reliacinė duomenų bazė.

Sistema patalpinta virtualiame serveryje, Digital Ocean debesų platformoje.
Klientai norėdami naudotis sistema, pirmiausia gaus serveryje talpinamą
statinį internetinės sąsajos failą - `index.html`. Faile esančios funkcijos
ir metodai komunikuos su aplikacijos serveriu naudojant axios paketetą. 
Browser ir backend procesai komunikuos HTTP aplikacijos lygio protokolu. 
Backend nenaudos ORM įrankių ir komunikuos su duomenu baze tiesiogiai, rašant
užklausas. Duomenų bazė bus paleista tame pačiame virtualiame serveryje, Docker
konteineryje.

![File](https://user-images.githubusercontent.com/58231312/190913623-01d242d5-cf43-4f82-9824-6cbe3586b6d3.jpg)

# API aprašymas

- `GET` `/version` - Versijos pasiimimas.
	- Reikia prisijungti: Ne
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija: Nėra
	- Atsakymo informacija:
	```JSON
	{
		"version": "1.0.0"
	}
	```

## Authorizacija

- `POST` `/login` - Prisijungimas.
	- Reikia prisijungti: Ne
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija:
	```JSON
	{
		"name": "name",
		"password": "password"
	}
	```
	- Atsakymo informacija:
	```JSON
	{
		"user": {
			"id": "cciuf5f6i1e0e49j5750",
			"name": "name",
			"admin": true,
			"created_at": "2022-09-15T19:53:42"
		},
		"access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImZvb2RpZSIsImV4cCI6Ijk1NTc5NjQ2MCIsImlhdCI6Ijk1NTc5Mjg2MCIsImlzcyI6ImZvb2RpZSIsInN1YiI6IjczczNyODc2aTFlNzJuNGgzZDQwIn0.0O2UJyLoKD3q-3DcbWfBvRhCY8OF00Ko462JUUHiaNU"
	}
	```

- `POST` `/register` - Registracija.
	- Reikia prisijungti: Ne
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija:
	```JSON
	{
		"name": "name",
		"password": "password"
	}
	```
	- Atsakymo informacija:
	```JSON
	{
		"user": {
			"id": "cciuf5f6i1e0e49j5750",
			"name": "name",
			"admin": true,
			"created_at": "2022-09-15T19:53:42"
		},
		"access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImZvb2RpZSIsImV4cCI6Ijk1NTc5NjQ2MCIsImlhdCI6Ijk1NTc5Mjg2MCIsImlzcyI6ImZvb2RpZSIsInN1YiI6IjczczNyODc2aTFlNzJuNGgzZDQwIn0.0O2UJyLoKD3q-3DcbWfBvRhCY8OF00Ko462JUUHiaNU"
	}
	```

- Endpoint'ai, kuriems yra reikalinga authorizacija (prisijungimas), turi turėti `Authorization` header. Pavyzdys:
`Authorization: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImZvb2RpZSIsImV4cCI6Ijk1NTc5NjQ2MCIsImlhdCI6Ijk1NTc5Mjg2MCIsImlzcyI6ImZvb2RpZSIsInN1YiI6IjczczNyODc2aTFlNzJuNGgzZDQwIn0.0O2UJyLoKD3q-3DcbWfBvRhCY8OF00Ko462JUUHiaNU"`

## Produktai

- `GET` `/products` - Produktų sąrašo pasiimimas.
	- Reikia prisijungti: Ne
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija: Nėra
	- Atsakymo informacija:
	```JSON
	[
		{
			"id": "cceqj5n6i1e7hgou9lv1",
			"name": "name",
			"serving": {
				"type": "grams",
				"size": 100,
				"calories": 350
			},
			"created_at": "2022-09-12T12:20:05"
		},
		{
			"id": "cceqj5n6i1e7hgou9lv2",
			"name": "name",
			"serving": {
				"type": "units",
				"size": 1,
				"calories": 150
			},
			"created_at": "2022-09-13T13:23:25"
		}
	]
	```

- `GET` `/products/{productID}` - Produkto pasiimimas.
	- Reikia prisijungti: Ne
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija: Nėra
	- Atsakymo informacija:
	```JSON
	{
		"id": "cceqj5n6i1e7hgou9lv1",
		"name": "name",
		"serving": {
			"type": "grams",
			"size": 100,
			"calories": 350
		},
		"created_at": "2022-09-12T12:20:05"
	}
	```

- `POST` `/products` - Produkto sukūrimas.
	- Reikia prisijungti: Taip
	- Reikalingos administratoriaus teisės: Taip
	- Užklausos informacija:
	```JSON
	{
		"name": "name",
		"serving": {
			"type": "grams",
			"size": 100,
			"calories": 350
		}
	}
	```
	- Atsakymo informacija:
	```JSON
	{
		"id": "cceqj5n6i1e7hgou9lv1",
		"name": "name",
		"serving": {
			"type": "grams",
			"size": 100,
			"calories": 350
		},
		"created_at": "2022-09-12T12:20:05"
	}
	```

- `PATCH` `/products/{productID}` - Produkto atnaujinimas.
	- Reikia prisijungti: Taip
	- Reikalingos administratoriaus teisės: Taip
	- Užklausos informacija:
	```JSON
	{
		"name": "name",
		"serving": {
			"type": "grams",
			"size": 100,
			"calories": 350
		}
	}
	```
	- Atsakymo informacija:
	```JSON
	{
		"id": "cceqj5n6i1e7hgou9lv1",
		"name": "name",
		"serving": {
			"type": "grams",
			"size": 100,
			"calories": 350
		},
		"created_at": "2022-09-12T12:20:05"
	}
	```

- `DELETE` `/products/{productID}` - Produkto ištrinimas.
	- Reikia prisijungti: Taip
	- Reikalingos administratoriaus teisės: Taip
	- Užklausos informacija: Nėra
	- Atsakymo informacija: Nėra

## Receptai

- `GET` `/recipes` - Receptų sąrašo pasiimimas.
	- Reikia prisijungti: Ne
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija: Nėra
	- Atsakymo informacija:
	```JSON
	[
		{
			"id": "cceqj5n6i1e7hgou9lv1",
			"user_id": "cceqj5n6i1e7hgou9lv2",
			"name": "name",
			"description": "description",
			"products": [
				{
					"product_id": "cceqj5n6i1e7hgou9lv0",
					"quantity": 2
				},
				{
					"product_id": "cceqj5n6i1e7hgou9lv9",
					"quantity": 1
				}
			],
			"created_at": "2022-09-12T12:20:05"
		},
		{
			"id": "cceqj5n6i1e7hgou9lv3",
			"user_id": "cceqj5n6i1e7hgou9lv4",
			"name": "name",
			"description": "description",
			"products": [
				{
					"product_id": "cceqj5n6i1e7hgou9lv5",
					"quantity": 2
				},
				{
					"product_id": "cceqj5n6i1e7hgou9lv6",
					"quantity": 1
				},
				{
					"product_id": "cceqj5n6i1e7hgou9lv0",
					"quantity": 5
				}
			],
			"created_at": "2022-09-12T18:50:25"
		}
	]
	```

- `GET` `/recipes/{recipyID}` - Recepto pasiimimas.
	- Reikia prisijungti: Ne
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija: Nėra
	- Atsakymo informacija:
	```JSON
	{
		"id": "cceqj5n6i1e7hgou9lv1",
		"user_id": "cceqj5n6i1e7hgou9lv2",
		"name": "name",
		"description": "description",
		"products": [
			{
				"product_id": "cceqj5n6i1e7hgou9lv0",
				"quantity": 2.5
			}
		],
		"created_at": "2022-09-12T12:20:05"
	}
	```

- `GET` `/recipes/user/{userID}` - Vartotojo receptų pasiimimas.
	- Reikia prisijungti: Ne
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija: Nėra
	- Atsakymo informacija:
	```JSON
	[
		{
			"id": "cceqj5n6i1e7hgou9lv1",
			"user_id": "cceqj5n6i1e7hgou9lv2",
			"name": "name",
			"description": "description",
			"products": [
				{
					"product_id": "cceqj5n6i1e7hgou9lv0",
					"quantity": 2
				},
				{
					"product_id": "cceqj5n6i1e7hgou9lv9",
					"quantity": 1
				}
			],
			"created_at": "2022-09-12T12:20:05"
		},
		{
			"id": "cceqj5n6i1e7hgou9lv3",
			"user_id": "cceqj5n6i1e7hgou9lv4",
			"name": "name",
			"description": "description",
			"products": [
				{
					"product_id": "cceqj5n6i1e7hgou9lv5",
					"quantity": 2
				},
				{
					"product_id": "cceqj5n6i1e7hgou9lv6",
					"quantity": 1
				},
				{
					"product_id": "cceqj5n6i1e7hgou9lv0",
					"quantity": 5
				}
			],
			"created_at": "2022-09-12T18:50:25"
		}
	]
	```

- `POST` `/recipes` - Recepto sukūrimas.
	- Reikia prisijungti: Taip
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija:
	```JSON
	{
		"name": "name",
		"description": "description",
		"products": [
			{
				"product_id": "cceqj5n6i1e7hgou9lv0",
				"quantity": 2.5
			}
		]
	}
	```
	- Atsakymo informacija:
	```JSON
	{
		"id": "cceqj5n6i1e7hgou9lv1",
		"user_id": "cceqj5n6i1e7hgou9lv2",
		"name": "name",
		"description": "description",
		"products": [
			{
				"product_id": "cceqj5n6i1e7hgou9lv0",
				"quantity": 2.5
			}
		],
		"created_at": "2022-09-12T12:20:05"
	}
	```

- `PATCH` `/recipes/{recipyID}` - Recepto atnaujinimas.
	- Reikia prisijungti: Taip
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija:
	```JSON
	{
		"name": "name",
		"description": "description",
		"products": [
			{
				"product_id": "cceqj5n6i1e7hgou9lv0",
				"quantity": 2.5
			}
		]
	}
	```
	- Atsakymo informacija:
	```JSON
	{
		"id": "cceqj5n6i1e7hgou9lv1",
		"user_id": "cceqj5n6i1e7hgou9lv2",
		"name": "name",
		"description": "description",
		"products": [
			{
				"product_id": "cceqj5n6i1e7hgou9lv0",
				"quantity": 2.5
			}
		],
		"created_at": "2022-09-12T12:20:05"
	}
	```

- `DELETE` `/recipes/{recipyID}` - Recepto ištrinimas.
	- Reikia prisijungti: Taip
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija: Nėra
	- Atsakymo informacija: Nėra

## Planai

- `GET` `/plans` - Planų sąrašo pasiimimas.
	- Reikia prisijungti: Ne
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija: Nėra
	- Atsakymo informacija:
	```JSON
	[
		{
			"id": "cceqj5n6i1e7hgou9lv1",
			"user_id": "cceqj5n6i1e7hgou9lv2",
			"name": "name",
			"description": "description",
			"recipes": [
				{
					"recipy_id": "cceqj5n6i1e7hgou9lv0",
					"quantity": 2
				},
				{
					"recipy_id": "cceqj5n6i1e7hgou9lv9",
					"quantity": 1
				}
			],
			"created_at": "2022-09-12T12:20:05"
		},
		{
			"id": "cceqj5n6i1e7hgou9lv3",
			"user_id": "cceqj5n6i1e7hgou9lv4",
			"name": "name",
			"description": "description",
			"recipes": [
				{
					"recipy_id": "cceqj5n6i1e7hgou9lv5",
					"quantity": 2
				},
				{
					"recipy_id": "cceqj5n6i1e7hgou9lv6",
					"quantity": 1
				},
				{
					"recipy_id": "cceqj5n6i1e7hgou9lv0",
					"quantity": 5
				}
			],
			"created_at": "2022-09-12T18:50:25"
		}
	]
	```

- `GET` `/plans/{planID}` - Plano pasiimimas.
	- Reikia prisijungti: Ne
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija: Nėra
	- Atsakymo informacija:
	```JSON
	{
		"id": "cceqj5n6i1e7hgou9lv1",
		"user_id": "cceqj5n6i1e7hgou9lv2",
		"name": "name",
		"description": "description",
		"recipes": [
			{
				"recipy_id": "cceqj5n6i1e7hgou9lv0",
				"quantity": 2
			},
			{
				"recipy_id": "cceqj5n6i1e7hgou9lv9",
				"quantity": 1
			}
		],
		"created_at": "2022-09-12T12:20:05"
	}
	```

- `GET` `/plans/user/{userID}` - Vartotojo planų pasiimimas.
	- Reikia prisijungti: Ne
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija: Nėra
	- Atsakymo informacija:
	```JSON
	[
		{
			"id": "cceqj5n6i1e7hgou9lv1",
			"user_id": "cceqj5n6i1e7hgou9lv2",
			"name": "name",
			"description": "description",
			"recipes": [
				{
					"recipy_id": "cceqj5n6i1e7hgou9lv0",
					"quantity": 2
				},
				{
					"recipy_id": "cceqj5n6i1e7hgou9lv9",
					"quantity": 1
				}
			],
			"created_at": "2022-09-12T12:20:05"
		},
		{
			"id": "cceqj5n6i1e7hgou9lv3",
			"user_id": "cceqj5n6i1e7hgou9lv4",
			"name": "name",
			"description": "description",
			"recipes": [
				{
					"recipy_id": "cceqj5n6i1e7hgou9lv5",
					"quantity": 2
				},
				{
					"recipy_id": "cceqj5n6i1e7hgou9lv6",
					"quantity": 1
				},
				{
					"recipy_id": "cceqj5n6i1e7hgou9lv0",
					"quantity": 5
				}
			],
			"created_at": "2022-09-12T18:50:25"
		}
	]
	```

- `POST` `/plans` - Plano sukūrimas.
	- Reikia prisijungti: Taip
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija:
	```JSON
	{
		"name": "name",
		"description": "description",
		"recipes": [
			{
				"recipy_id": "cceqj5n6i1e7hgou9lv0",
				"quantity": 2
			}
		]
	}
	```
	- Atsakymo informacija:
	```JSON
	{
		"id": "cceqj5n6i1e7hgou9lv1",
		"user_id": "cceqj5n6i1e7hgou9lv2",
		"name": "name",
		"description": "description",
		"recipes": [
			{
				"recipy_id": "cceqj5n6i1e7hgou9lv0",
				"quantity": 2
			}
		],
		"created_at": "2022-09-12T12:20:05"
	}
	```

- `PATCH` `/plans/{planID}` - Plano atnaujinimas.
	- Reikia prisijungti: Taip
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija:
	```JSON
	{
		"name": "name",
		"description": "description",
		"recipes": [
			{
				"recipy_id": "cceqj5n6i1e7hgou9lv0",
				"quantity": 2
			}
		]
	}
	```
	- Atsakymo informacija:
	```JSON
	{
		"id": "cceqj5n6i1e7hgou9lv1",
		"user_id": "cceqj5n6i1e7hgou9lv2",
		"name": "name",
		"description": "description",
		"recipes": [
			{
				"recipy_id": "cceqj5n6i1e7hgou9lv0",
				"quantity": 2
			}
		],
		"created_at": "2022-09-12T12:20:05"
	}
	```

- `DELETE` `/plans/{planID}` - Plano ištrinimas.
	- Reikia prisijungti: Taip
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija: Nėra
	- Atsakymo informacija: Nėra

## Vartotojai

- `DELETE` `/users` - Savo vartotojo ištrinimas.
	- Reikia prisijungti: Taip
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija: Nėra
	- Atsakymo informacija: Nėra

- `PATCH` `/users` - Vartotojo slaptažodžio atnaujinimas.
	- Reikia prisijungti: Taip
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija:
	```JSON
	{
		"password": "password",
		"old_password": "old_password"
	}
	```

Atsakymo informacija: Nėra

- `GET` `/users` - Vartotojų pasiimimas.
	- Reikia prisijungti: Taip
	- Reikalingos administratoriaus teisės: Taip
	- Užklausos informacija: Nėra
	- Atsakymo informacija:
	```JSON
	[
		{
			"id": "cciuf5f6i1e0e49j5750",
			"name": "name",
			"admin": true,
			"created_at": "2022-09-15T19:53:42"
		},
		{
			"id": "cciuf5f6i1e0e49j5752",
			"name": "name2",
			"admin": false,
			"created_at": "2022-09-16T21:22:12"
		}
	]
	```

- `POST` `/users` - Administratoriaus vartotojo sukūrimas.
	- Reikia prisijungti: Taip
	- Reikalingos administratoriaus teisės: Taip
	- Užklausos informacija: 
	```JSON
	{
		"name": "name",
		"password": "password"
	}
	```
	- Atsakymo informacija:
	```JSON
	{
		"id": "cciuf5f6i1e0e49j5750",
		"name": "name",
		"admin": true,
		"created_at": "2022-09-15T19:53:42"
	}
	```

- `GET` `/users/{userID}` - Vartotojaus pasiimimas.
	- Reikia prisijungti: Taip
	- Reikalingos administratoriaus teisės: Taip
	- Užklausos informacija: Nėra
	- Atsakymo informacija:
	```JSON
	{
		"id": "cciuf5f6i1e0e49j5750",
		"name": "name",
		"admin": true,
		"created_at": "2022-09-15T19:53:42"
	}
	```

- `DELETE` `/users/{userID}` - Vartotojaus ar administratoriaus ištrinimas.
	- Reikia prisijungti: Taip
	- Reikalingos administratoriaus teisės: Taip
	- Užklausos informacija: Nėra
	- Atsakymo informacija: Nėra
