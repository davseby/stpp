import { ref } from "vue"
import { defineStore } from "pinia"
import { useAuth } from "@/stores/useAuth"
import { useProduct } from "@/stores/useProduct"
import axios from 'axios'

export const useRecipe = defineStore("recipe", () => {
	interface Recipe {
		id: string
		name: string
		description: string
		image_url: string
		products: {
			product_id: string
			quantity: number
		}
	}

	const recipes = ref<Array<Recipe>>([]);

	function retrieveRecipes() {
		return useProduct().retrieveProducts().then((response) => {
			if (response.error) {
				return response
			}

			let products = response.data;

			return axios
				.get("/api/recipes")
				.then((response) => {
					if (response.status == 200) {
						let new_recipes = response.data;
						for (let i = 0; i < new_recipes.length; i++) {
							for (let j = 0; j < new_recipes[i].products.length; j++) {
								for (let k = 0; k < products.length; k++) {
									if (products[k].id == new_recipes[i].products[j].product_id) {
										const recipeProduct = new_recipes[i].products[j];
										new_recipes[i].products[j] = JSON.parse(JSON.stringify(products[k]));
										new_recipes[i].products[j].quantity = recipeProduct.quantity;
									}
								}
							}
						}

						recipes.value = new_recipes;
					}

					return {
						data: response.data,
					}
				})
				.catch((error) => {
					return {
						error: error.message,
					}
				})
		}).catch((error) => {
			return {
				error: error.message,
			}
		});

	}

	function createRecipe(recipe: Recipe) {
		return axios
			.post("/api/recipes", recipe, {
				headers: {
					Authorization: 'Bearer ' + useAuth().key,
				}
			})
			.then((response) => {
				if (response.status == 200) {
					retrieveRecipes();
				}

				return {
					data: response.data,
				}
			})
			.catch((error) => {
				return {
					error: error.message,
				}
			})
	}

	function updateRecipe(recipe: Recipe) {
		return axios
			.patch("/api/recipes/" + recipe.id, recipe, {
				headers: {
					Authorization: 'Bearer ' + useAuth().key,
				}
			})
			.then((response) => {
				if (response.status == 200) {
					retrieveRecipes();
				} 

				return {
					data: response.data,
				}
			})
			.catch((error) => {
				return {
					error: error.message,
				}
			})
	}

	function deleteRecipe(id: string) {
		return axios
			.delete("/api/recipes/" + id, {
				headers: {
					Authorization: 'Bearer ' + useAuth().key,
				}
			})
			.then((response) => {
				if (response.status == 204) {
					retrieveRecipes();
				}

				return {}
			})
			.catch((error) => {
				if (error.response.status == 409) {
					return {
						error: "Recipe is in use.",
					}
				}

				return {
					error: error.message,
				}
			})
	}

	return { recipes, retrieveRecipes, createRecipe, updateRecipe, deleteRecipe }
})
