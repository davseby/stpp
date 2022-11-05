import { ref } from "vue"
import { defineStore } from "pinia"
import { useAuth } from "@/stores/useAuth"
import { useRecipe } from "@/stores/useRecipe"
import axios from 'axios'

export const usePlan = defineStore("plan", () => {
	interface Plan {
		id: string
		name: string
		description: string
		image_url: string
		recipes: {
			recipe_id: string
			quantity: number
		}
	}

	const plans = ref<Array<Plan>>([]);

	function retrievePlans() {
		return useRecipe().retrieveRecipes().then((response) => {
			if (response.error) {
				return response
			}

			let recipes = response.data;

			return axios
				.get("/api/plans")
				.then((response) => {
					if (response.status == 200) {
						let new_plans = response.data;
						for (let i = 0; i < new_plans.length; i++) {
							for (let j = 0; j < new_plans[i].recipes.length; j++) {
								for (let k = 0; k < recipes.length; k++) {
									if (recipes[k].id == new_plans[i].recipes[j].recipe_id) {
										const planRecipe = new_plans[i].recipes[j];
										new_plans[i].recipes[j] = JSON.parse(JSON.stringify(recipes[k]));
										new_plans[i].recipes[j].quantity = planRecipe.quantity;
									}
								}
							}
						}

						plans.value = new_plans;
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

	function createPlan(plan: Plan) {
		return axios
			.post("/api/plans", plan, {
				headers: {
					Authorization: 'Bearer ' + useAuth().key,
				}
			})
			.then((response) => {
				if (response.status == 200) {
					retrievePlans();
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

	function updatePlan(plan: Plan) {
		return axios
			.patch("/api/plans/" + plan.id, plan, {
				headers: {
					Authorization: 'Bearer ' + useAuth().key,
				}
			})
			.then((response) => {
				if (response.status == 200) {
					retrievePlans();
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

	function deletePlan(id: string) {
		return axios
			.delete("/api/plans/" + id, {
				headers: {
					Authorization: 'Bearer ' + useAuth().key,
				}
			})
			.then((response) => {
				if (response.status == 204) {
					retrievePlans();
				}

				return {}
			})
			.catch((error) => {
				return {
					error: error.message,
				}
			})
	}

	return { plans, retrievePlans, createPlan, updatePlan, deletePlan }
})
