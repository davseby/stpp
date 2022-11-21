Modulio kodas: T120B165

Dėstytojai:
- Doc. Prakt. Rasa Mažutienė
- Doc. Prakt. Petras Tamošiūnas

Studentas:
- Dovydas Bykovas, IFF-9/2 grupė

# Foodie (Saityno taikomųjų programų projektavimas)

Projekto tikslas - viešą maisto receptų bei jų planų kūrimo ir 
dalijimosi platformą.

Projekto veikimo principas - administratoriai kuria produktus iš kurių vartotojai
kuria receptus, o iš receptų - maisto planus. Visa informacija yra viešai prieinama.

Programiškas projekto veikimo principas - platformą sudaro internetinė sąsaja 
(frontend) ir serverio aplikacija (backend). Internetinės sąsajos ir serverio 
aplikacija komunikacija vyksta per RESTful API.

Dizaino failas - https://github.com/davseby/stpp/blob/master/DESIGN.md

Wireframe failas - https://github.com/davseby/stpp/blob/master/WIREFRAME.md

Sistema turi tris roles:

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
- Atnaujinti savo receptą.
- Ištrinti savo receptą.
- Sukurti planą.
- Atnaujinti savo planą.
- Ištrinti savo planą.
- Ištrinti savo vartotoją.
- Pasikeisti savo slaptažodį.
- Pasiimti vartotojo informaciją.

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

Papildomos funkcijos:
- Saugomas slaptažodis turi būti užšifruotas.
- Negalima ištrinti pagrindinio administratoriaus.
- Vartotojas negali ištrinti produkto ar recepto, jeigu jis yra naudojamas.
- Jeigu yra ištrinamas vartotojas, jo receptai ir planai nenurodo sukurėjo identifikacijos.

# Sistemos architektūra

Aplikacijos frontend dalį sudaro TypeScript ir Vue 3 karkasas.
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

![File2](https://user-images.githubusercontent.com/58231312/190964890-c02acb5e-27cc-4dcf-af2a-aa9bfcaea533.png)


# API aprašymas

- `GET` `/api/version` - Versijos pasiimimas.
	- Reikia prisijungti: Ne
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija: Nėra
	- Galimi atsakymai:
		- `200`
		```JSON
		{
			"version": "1.0.0"
		}
		```

## Autorizacija ir autentikacija

- `POST` `/api/login` - Prisijungimas.
	- Reikia prisijungti: Ne
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija:
	```JSON
	{
		"name": "name",
		"password": "password"
	}
	```
	- Galimi atsakymai:
		- `200` prijungtas vartotojas.
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
		- `400` netinkama informacija.
		- `401` nerastas vartotojas.
		- `500` serverio klaida.

- `POST` `/api/register` - Registracija.
	- Reikia prisijungti: Ne
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija:
	```JSON
	{
		"name": "name",
		"password": "password"
	}
	```
	- Galimi atsakymai:
		- `200` užregistruotas vartotojas.
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
		- `400` netinkama informacija.
		- `409` vartotojas su tokiu vardu jau egzistuoja.
		- `500` serverio klaida.

- Endpoint'ai, kuriems yra reikalinga authorizacija (prisijungimas), turi turėti `Authorization` header. Pavyzdys:
```
Authorization: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImZvb2RpZSIsImV4cCI6Ijk1NTc5NjQ2MCIsImlhdCI6Ijk1NTc5Mjg2MCIsImlzcyI6ImZvb2RpZSIsInN1YiI6IjczczNyODc2aTFlNzJuNGgzZDQwIn0.0O2UJyLoKD3q-3DcbWfBvRhCY8OF00Ko462JUUHiaNU"
```

## Produktai

- `GET` `/api/products` - Produktų sąrašo pasiimimas.
	- Reikia prisijungti: Ne
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija: Nėra
	- Galimi atsakymai:
		- `200` produktų sąrašas.
		```JSON
		[
			{
				"id": "cceqj5n6i1e7hgou9lv1",
				"name": "name",
				"image_url": "image_url",
				"description": "description",
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
				"image_url": "image_url",
				"description": "description",
				"serving": {
					"type": "units",
					"size": 1,
					"calories": 150
				},
				"created_at": "2022-09-13T13:23:25"
			}
		]
		```
		- `500` serverio klaida.

- `GET` `/api/products/{productID}` - Produkto pasiimimas.
	- Reikia prisijungti: Ne
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija: Nėra
	- Galimi atsakymai:
		- `200` vienas produktas.
		```JSON
		{
			"id": "cceqj5n6i1e7hgou9lv1",
			"name": "name",
			"image_url": "image_url",
			"description": "description",
			"serving": {
				"type": "grams",
				"size": 100,
				"calories": 350
			},
			"created_at": "2022-09-12T12:20:05"
		}
		```
		- `400` blogas produkto id.
		- `404` produktas neegzistuoja.
		- `500` serverio klaida.

- `POST` `/api/products` - Produkto sukūrimas.
	- Reikia prisijungti: Taip
	- Reikalingos administratoriaus teisės: Taip
	- Užklausos informacija:
	```JSON
	{
		"name": "name",
		"image_url": "image_url",
		"description": "description",
		"serving": {
			"type": "grams",
			"size": 100,
			"calories": 350
		}
	}
	```
	- Galimi atsakymai:
		- `200` sukurtas produktas.
		```JSON
		{
			"id": "cceqj5n6i1e7hgou9lv1",
			"name": "name",
			"image_url": "image_url",
			"description": "description",
			"serving": {
				"type": "grams",
				"size": 100,
				"calories": 350
			},
			"created_at": "2022-09-12T12:20:05"
		}
		```
		- `400` bloga produkto informacija.
		- `500` serverio klaida.

- `PATCH` `/api/products/{productID}` - Produkto atnaujinimas.
	- Reikia prisijungti: Taip
	- Reikalingos administratoriaus teisės: Taip
	- Užklausos informacija:
	```JSON
	{
		"name": "name",
		"image_url": "image_url",
		"description": "description",
		"serving": {
			"type": "grams",
			"size": 100,
			"calories": 350
		}
	}
	```
	- Galimi atsakymai:
		- `200` atnaujintas produktas.
		```JSON
		{
			"id": "cceqj5n6i1e7hgou9lv1",
			"name": "name",
			"image_url": "image_url",
			"description": "description",
			"serving": {
				"type": "grams",
				"size": 100,
				"calories": 350
			},
			"created_at": "2022-09-12T12:20:05"
		}
		```
		- `400` bloga produkto informacija arba blogas id.
		- `500` serverio klaida.

- `DELETE` `/api/products/{productID}` - Produkto ištrinimas.
	- Reikia prisijungti: Taip
	- Reikalingos administratoriaus teisės: Taip
	- Užklausos informacija: Nėra
	- Galimi atsakymai:
		- `204` ištrintas produktas.
		- `400` blogas produkto id.
		- `404` nerastas produktas.
		- `409` naudojamas produktas.
		- `500` serverio klaida.

## Receptai

- `GET` `/api/recipes` - Receptų sąrašo pasiimimas.
	- Reikia prisijungti: Ne
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija: Nėra
	- Galimi atsakymai:
		- `200` receptų sąrašas.
		```JSON
		[
			{
				"id": "cceqj5n6i1e7hgou9lv1",
				"user_id": "cceqj5n6i1e7hgou9lv2",
				"name": "name",
				"image_url": "image_url",
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
				"image_url": "image_url",
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
		- `500` serverio klaida.

- `GET` `/api/recipes/{recipeID}` - Recepto pasiimimas.
	- Reikia prisijungti: Ne
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija: Nėra
	- Galimi atsakymai:
		- `200` vienas receptas.
		```JSON
		{
			"id": "cceqj5n6i1e7hgou9lv1",
			"user_id": "cceqj5n6i1e7hgou9lv2",
			"name": "name",
			"image_url": "image_url",
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
		- `400` blogas recepto id.
		- `404` receptas neegzistuoja.
		- `500` serverio klaida.

- `POST` `/api/recipes` - Recepto sukūrimas.
	- Reikia prisijungti: Taip
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija:
	```JSON
	{
		"name": "name",
		"image_url": "image_url",
		"description": "description",
		"products": [
			{
				"product_id": "cceqj5n6i1e7hgou9lv0",
				"quantity": 2.5
			}
		]
	}
	```
	- Galimi atsakymai:
		- `200` sukurtas receptas.
		```JSON
		{
			"id": "cceqj5n6i1e7hgou9lv1",
			"user_id": "cceqj5n6i1e7hgou9lv2",
			"name": "name",
			"image_url": "image_url",
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
		- `400` bloga recepto informacija.
		- `500` serverio klaida.

- `PATCH` `/api/recipes/{recipeID}` - Recepto atnaujinimas.
	- Reikia prisijungti: Taip
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija:
	```JSON
	{
		"name": "name",
		"image_url": "image_url",
		"description": "description",
		"products": [
			{
				"product_id": "cceqj5n6i1e7hgou9lv0",
				"quantity": 2.5
			}
		]
	}
	```
	- Galimi atsakymai:
		- `200` atnaujintas receptas.
		```JSON
		{
			"id": "cceqj5n6i1e7hgou9lv1",
			"user_id": "cceqj5n6i1e7hgou9lv2",
			"name": "name",
			"image_url": "image_url",
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
		- `400` bloga recepto informacija arba blogas id.
		- `500` serverio klaida.

- `DELETE` `/api/recipes/{recipeID}` - Recepto ištrinimas.
	- Reikia prisijungti: Taip
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija: Nėra
	- Galimi atsakymai:
		- `204` ištrintas receptas.
		- `400` blogas recepto id.
		- `404` nerastas receptas.
		- `409` naudojamas receptas.
		- `500` serverio klaida.

## Planai

- `GET` `/api/plans` - Planų sąrašo pasiimimas.
	- Reikia prisijungti: Ne
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija: Nėra
	- Galimi atsakymai:
		- `200` planų sąrašas.
		```JSON
		[
			{
				"id": "cceqj5n6i1e7hgou9lv1",
				"user_id": "cceqj5n6i1e7hgou9lv2",
				"name": "name",
				"description": "description",
				"recipes": [
					{
						"recipe_id": "cceqj5n6i1e7hgou9lv0",
						"quantity": 2
					},
					{
						"recipe_id": "cceqj5n6i1e7hgou9lv9",
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
						"recipe_id": "cceqj5n6i1e7hgou9lv5",
						"quantity": 2
					},
					{
						"recipe_id": "cceqj5n6i1e7hgou9lv6",
						"quantity": 1
					},
					{
						"recipe_id": "cceqj5n6i1e7hgou9lv0",
						"quantity": 5
					}
				],
				"created_at": "2022-09-12T18:50:25"
			}
		]
		```
		- `500` serverio klaida.

- `GET` `/api/plans/{planID}` - Plano pasiimimas.
	- Reikia prisijungti: Ne
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija: Nėra
	- Galimi atsakymai:
		- `200` vienas planas.
		```JSON
		{
			"id": "cceqj5n6i1e7hgou9lv1",
			"user_id": "cceqj5n6i1e7hgou9lv2",
			"name": "name",
			"description": "description",
			"recipes": [
				{
					"recipe_id": "cceqj5n6i1e7hgou9lv0",
					"quantity": 2
				},
				{
					"recipe_id": "cceqj5n6i1e7hgou9lv9",
					"quantity": 1
				}
			],
			"created_at": "2022-09-12T12:20:05"
		}
		```
		- `400` blogas plano id.
		- `404` planas neegzistuoja.
		- `500` serverio klaida.

- `POST` `/api/plans` - Plano sukūrimas.
	- Reikia prisijungti: Taip
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija:
	```JSON
	{
		"name": "name",
		"description": "description",
		"recipes": [
			{
				"recipe_id": "cceqj5n6i1e7hgou9lv0",
				"quantity": 2
			}
		]
	}
	```
	- Galimi atsakymai:
		- `200` sukurtas planas.
		```JSON
		{
			"id": "cceqj5n6i1e7hgou9lv1",
			"user_id": "cceqj5n6i1e7hgou9lv2",
			"name": "name",
			"description": "description",
			"recipes": [
				{
					"recipe_id": "cceqj5n6i1e7hgou9lv0",
					"quantity": 2
				}
			],
			"created_at": "2022-09-12T12:20:05"
		}
		```
		- `400` bloga plano informacija.
		- `500` serverio klaida.

- `PATCH` `/api/plans/{planID}` - Plano atnaujinimas.
	- Reikia prisijungti: Taip
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija:
	```JSON
	{
		"name": "name",
		"description": "description",
		"recipes": [
			{
				"recipe_id": "cceqj5n6i1e7hgou9lv0",
				"quantity": 2
			}
		]
	}
	```
	- Galimi atsakymai:
		- `200` atnaujintas planas.
		```JSON
		{
			"id": "cceqj5n6i1e7hgou9lv1",
			"user_id": "cceqj5n6i1e7hgou9lv2",
			"name": "name",
			"description": "description",
			"recipes": [
				{
					"recipe_id": "cceqj5n6i1e7hgou9lv0",
					"quantity": 2
				}
			],
			"created_at": "2022-09-12T12:20:05"
		}
		```
		- `400` bloga plano informacija arba blogas id.
		- `500` serverio klaida.

- `DELETE` `/api/plans/{planID}` - Plano ištrinimas.
	- Reikia prisijungti: Taip
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija: Nėra
	- Galimi atsakymai:
		- `204` ištrintas planas.
		- `400` blogas plano id.
		- `404` nerastas planas.
		- `409` naudojamas planas.
		- `500` serverio klaida.

## Vartotojai

- `DELETE` `/api/users` - Savo vartotojo ištrinimas.
	- Reikia prisijungti: Taip
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija: Nėra
	- Galimi atsakymai:
		- `204` ištrintas vartotojas.
		- `500` serverio klaida.

- `PATCH` `/api/users` - Vartotojo slaptažodžio atnaujinimas.
	- Reikia prisijungti: Taip
	- Reikalingos administratoriaus teisės: Ne
	- Užklausos informacija:
	```JSON
	{
		"password": "password",
		"old_password": "old_password"
	}
	```
	- Galimi atsakymai:
		- `204` pakeistas vartotojo slaptažodis.
		- `400` nesutampantis senas arba netinkamas naujas slaptažodis.
		- `500` serverio klaida.

- `GET` `/api/users` - Vartotojų pasiimimas.
	- Reikia prisijungti: Taip
	- Reikalingos administratoriaus teisės: Taip
	- Užklausos informacija: Nėra
	- Galimi atsakymai:
		- `200` vartotojų sąrašas.
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
		- `500` serverio klaida.

- `POST` `/api/users` - Administratoriaus vartotojo sukūrimas.
	- Reikia prisijungti: Taip
	- Reikalingos administratoriaus teisės: Taip
	- Užklausos informacija: 
	```JSON
	{
		"name": "name",
		"password": "password"
	}
	```
	- Galimi atsakymai:
		- `200` sukurtas vartotojas.
		```JSON
		{
			"id": "cciuf5f6i1e0e49j5750",
			"name": "name",
			"admin": true,
			"created_at": "2022-09-15T19:53:42"
		}
		```
		- `400` bloga informacija.
		- `409` egzistuojantis vartotojo vardas.
		- `500` serverio klaida.

- `GET` `/api/users/{userID}` - Vartotojaus pasiimimas.
	- Reikia prisijungti: Taip
	- Reikalingos administratoriaus teisės: Taip
	- Užklausos informacija: Nėra
	- Galimi atsakymai:
		- `200` vienas vartotojas.
		```JSON
		{
			"id": "cciuf5f6i1e0e49j5750",
			"name": "name",
			"admin": true,
			"created_at": "2022-09-15T19:53:42"
		}
		```
		- `400` blogas id.
		- `404` nerastas vartotojas.
		- `500` serverio klaida.

- `DELETE` `/api/users/{userID}` - Vartotojaus ar administratoriaus ištrinimas.
	- Reikia prisijungti: Taip
	- Reikalingos administratoriaus teisės: Taip
	- Užklausos informacija: Nėra
	- Galimi atsakymai:
		- `204` ištrintas vartotojas.
		- `400` blogas id arba trinamas paskutinis administracinis vartotojas.
		- `404` nerastas vartotojas.
		- `500` serverio klaida.

# Išvados

Sistema pavyko įgyvendinti naudojant Go 1.19, TypeScript and Vue3 karkasu. 
Duomenų bazei buvo panaudota MariaDB, o lengvesniam development pasitelkta
Docker įrankis. Serveris ir statiniai jo failai buvo paleisti DigitealOcean
debesų kompiuterijos sistemoje ir buvo viešai prieinami internetu. Projekte
įgyvendintas visas užsibrėžtas funkcionalumas.
