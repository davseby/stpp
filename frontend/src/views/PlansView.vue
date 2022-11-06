<script setup lang="ts">
import { ref, computed } from 'vue'
import { useAuth } from "@/stores/useAuth"
import { useRecipe } from "@/stores/useRecipe"
import { usePlan } from "@/stores/usePlan"

import {
	useNotification,
	NForm,
	NFormItem,
	NImage,
	NInput,
	NPopconfirm,
	NInputNumber,
	NEllipsis,
	NSelect,
	NSpace,
	NButton,
	NIcon,
	NGrid,
	NGridItem,
	NCheckbox,
	NCard,
	NModal,
} from "naive-ui"

import {
	AddCircleOutline as CreateIcon,
	TrashOutline as DeleteIcon,
	SettingsOutline as UpdateIcon,
	EyeOutline as ViewIcon,
} from '@vicons/ionicons5'

const createPlanFormRef = ref<Object>(null);
const updatePlanFormRef = ref<Object>(null);

const planRules = {
	name: {
		message: "Please input a plan name",
		required: true,
	},
	recipes: {
		message: "Please select at least 1 recipe",
		validator(rule, value, callback) {
			return value.length > 0;
		},
	},
	description: {
		message: "Please describe the plan",
		required: true,
	},
}

const availableRecipes = ref<Array<Object>>([]);

const createdPlan = ref<Plan>({});
const createdPlanRecipes = ref<Array<Object>>([]);

const selectedPlan = ref<Plan>({});
const selectedPlanRecipes = ref<Array<Object>>([]);

const createModalVisible = ref(false);
const updateModalVisible = ref(false);

const notification = useNotification();

usePlan().retrievePlans();

const deletePlan = (id: string) => {
	usePlan().deletePlan(id).then(resp => {
		if (resp.error) {
			notification["error"]({
				content: resp.error,
				duration: 2500,
				keepAliveOnHover: true
			})
			return;
		} 

		notification["success"]({
			content: "Succesfully deleted a plan",
			duration: 2500,
			keepAliveOnHover: true
		})
	});
}

const showCreateModal = () => {
	createdPlan.value = {
		recipes: [],
	};

	createdPlanRecipes.value = [];
	for (let i = 0; i < createdPlan.value.recipes.length; i++) {
		createdPlanRecipes.value.push(createdPlan.value.recipes[i].id);
	}

	availableRecipes.value = [];
	for (let i = 0; i < useRecipe().recipes.length; i++) {
		availableRecipes.value.push({
			label: useRecipe().recipes[i].name,
			value: useRecipe().recipes[i].id,
		});
	}

	createModalVisible.value = true;
}

const viewPlan = ref(true);

const showUpdateModal = (plan: Plan, view: boolean) => {
	selectedPlan.value = JSON.parse(JSON.stringify(plan));

	selectedPlanRecipes.value = [];
	for (let i = 0; i < selectedPlan.value.recipes.length; i++) {
		selectedPlanRecipes.value.push(selectedPlan.value.recipes[i].id);
	}

	availableRecipes.value = [];
	for (let i = 0; i < useRecipe().recipes.length; i++) {
		availableRecipes.value.push({
			label: useRecipe().recipes[i].name,
			value: useRecipe().recipes[i].id,
		});
	}

	viewPlan.value = view;
	updateModalVisible.value = true;
}

const updateCreatedPlanRecipes = (value: string, option: any) => {
	let recipes = [];
	for (let i = 0; i < option.length; i++) {
		let found = false;
		for (let j = 0; j < createdPlan.value.recipes.length; j++) {
			if (option[i].value == createdPlan.value.recipes[j].id) {
				recipes.push(createdPlan.value.recipes[j]);
				found = true;
				break;
			}
		}

		if (!found) {
			for(let j = 0; j < useRecipe().recipes.length; j++) {
				if(option[i].value == useRecipe().recipes[j].id) {
					let recipe = JSON.parse(JSON.stringify(useRecipe().recipes[j]));
					recipe.quantity = 1;
					recipes.push(recipe);
					break;
				}
			}
		}
	}

	createdPlan.value.recipes = recipes;
}

const updateSelectedPlanRecipes = (value: string, option: any) => {
	let recipes = [];
	for (let i = 0; i < option.length; i++) {
		let found = false;
		for (let j = 0; j < selectedPlan.value.recipes.length; j++) {
			if (option[i].value == selectedPlan.value.recipes[j].id) {
				recipes.push(selectedPlan.value.recipes[j]);
				found = true;
				break;
			}
		}

		if (!found) {
			for(let j = 0; j < useRecipe().recipes.length; j++) {
				if(option[i].value == useRecipe().recipes[j].id) {
					let recipe = JSON.parse(JSON.stringify(useRecipe().recipes[j]));
					recipe.quantity = 1;
					recipes.push(recipe);
					break;
				}
			}
		}
	}

	selectedPlan.value.recipes = recipes;
}

const createPlan = () => {
	createPlanFormRef.value.validate((errors) => {
		if (errors) {
			return
		}

		let recipes = [];
		for (let i = 0; i < createdPlan.value.recipes.length; i++) {
			recipes.push({
				recipe_id: createdPlan.value.recipes[i].id,
				quantity: createdPlan.value.recipes[i].quantity,
			})
		}

		let plan = JSON.parse(JSON.stringify(createdPlan.value));
		plan.recipes = recipes;
		usePlan().createPlan(plan).then(resp => {
			if (resp.error) {
				notification["error"]({
					content: resp.error,
					duration: 2500,
					keepAliveOnHover: true
				})
				return;
			} 

			notification["success"]({
				content: "Succesfully created a plan",
				duration: 2500,
				keepAliveOnHover: true
			})

			createModalVisible.value = false;
		});
	}).catch((errors) => {
		console.log(errors);
	})
}

const updatePlan = () => {
	updatePlanFormRef.value.validate((errors) => {
		if (errors) {
			return
		}

		let recipes = [];
		for (let i = 0; i < selectedPlan.value.recipes.length; i++) {
			recipes.push({
				recipe_id: selectedPlan.value.recipes[i].id,
				quantity: selectedPlan.value.recipes[i].quantity,
			})
		}

		let plan = JSON.parse(JSON.stringify(selectedPlan.value));
		plan.recipes = recipes;

		usePlan().updatePlan(plan).then(resp => {
			if (resp.error) {
				notification["error"]({
					content: resp.error,
					duration: 2500,
					keepAliveOnHover: true
				})
				return;
			} 

			notification["success"]({
				content: "Succesfully updated a plan",
				duration: 2500,
				keepAliveOnHover: true
			})

			updateModalVisible.value = false;
		});
	}).catch((errors) => {
		console.log(errors);
	})
}

const searchQuery = ref("");
const searchSelfCreated = ref<boolean>(false);

const searchedPlans = computed(() => {
	return usePlan().plans.filter((plan) => {
		return (
			(!searchSelfCreated.value || plan.user_id == useAuth().user.id) && plan.name
				.toLowerCase()
				.indexOf(searchQuery.value.toLowerCase()) != -1
		);
	});
});

</script>

<template>
	<n-space justify="space-between" align="center">
		<div class="page-title">Plans</div>
		<n-space justify="space-between" align="center">
			<n-checkbox v-if="useAuth().authorized" v-model:checked="searchSelfCreated">Self Created</n-checkbox>
			<n-input v-model:value="searchQuery" type="text" placeholder="Search" />
			<n-button v-if="useAuth().authorized" :focusable="false" @click="showCreateModal" text>
				<template #icon>
					<n-icon size="20px" color="#36ad6a">
						<create-icon />
					</n-icon>
				</template>
				Create
			</n-button>
		</n-space>
	</n-space>
	<n-grid cols="3 800:4 1000:5 1200:6" x-gap="20px" y-gap="20px">
		<n-grid-item v-for="(plan, index) in searchedPlans" :key="index" style="display: flex">
			<n-card :title="plan.name" hoverable embedded>
				<n-ellipsis line-clamp="3" :tooltip="{width:200}">
					{{ plan.description }}
				</n-ellipsis>
				<div>
					<n-ellipsis line-clamp="5" :tooltip="{width:200}">
						<div v-for="(recipe, index) in plan.recipes" :key="index"> 
							<div class="nutrition-header">{{ recipe.quantity }} {{ recipe.name }}</div>
						</div>
					</n-ellipsis>
				</div>
				<template #footer>
					<n-space justify="space-between" align="center">
						<n-button :focusable="false" @click="showUpdateModal(plan, true)" text>
							<template #icon>
								<n-icon size="20px" color="#36ad6a">
									<view-icon />
								</n-icon>
							</template>
							View
						</n-button>
						<n-button v-if="useAuth().authorized && useAuth().user.id == plan.user_id" :focusable="false" @click="showUpdateModal(plan, false)" text>
							<template #icon>
								<n-icon size="20px" color="#36ad6a">
									<update-icon />
								</n-icon>
							</template>
							Update
						</n-button>
						<div v-else></div>
						<n-popconfirm v-if="useAuth().authorized && (useAuth().user.admin || useAuth().user.id == plan.user_id)"  width="310px" placement="top-end" :positive-text="null" negative-text="Confirm" @negative-click="deletePlan(plan.id)">
							<template #trigger>
								<n-button type="error" :focusable="false" text>
									<template #icon>
										<n-icon size="20px" color="#d03050">
											<delete-icon />
										</n-icon>
									</template>
									Delete
								</n-button>
							</template>
							Are you sure you want to delete the plan?
						</n-popconfirm>
					</n-space>
				</template>
			</n-card>
		</n-grid-item>
	</n-grid>
	<n-modal
		v-model:show="createModalVisible" 
		preset="dialog" 
		:show-icon="false"
	>
		<div class="modal-title">
			Create a new plan
		</div>
		<div class="modal-content">
			<n-form ref="createPlanFormRef" :model="createdPlan" :rules="planRules">
				<n-form-item label="Name" path="name" required>
					<n-input type="text" maxlength="255" show-count v-model:value="createdPlan.name" placeholder="Name of the plan" />
				</n-form-item>
				<n-form-item label="Recipes" path="recipes" required>
					<n-select max-tag-count="responsive" multiple v-model:value="createdPlanRecipes" @update:value="updateCreatedPlanRecipes" filterable :options="availableRecipes" />
				</n-form-item>
				<n-form-item v-if="createdPlanRecipes.length > 0" label="Recipes Quantities" required>
					<n-space vertical style="width: 100%; max-height: 300px; overflow-y: auto">
						<n-form-item v-for="(recipe, index) in createdPlan.recipes" :key="index" :label="recipe.name" :show-feedback="false" required>
							<n-input-number :show-buttom="false" :default-value="1" min=1 :precision="0" step=1 v-model:value="recipe.quantity" style="width: 100%" />
						</n-form-item>
					</n-space>
				</n-form-item>
				<n-form-item label="Description" path="description">
					<n-input type="textarea" maxlength="1023" show-count v-model:value="createdPlan.description" placeholder="Description of the plan" />
				</n-form-item>
			</n-form>
		</div>
		<template #action>
			<n-button type="primary" @click="createPlan">
				Create
			</n-button>
		</template>
	</n-modal>
	<n-modal
		v-model:show="updateModalVisible" 
		preset="dialog" 
		:show-icon="false"
	>
		<div v-if="!viewPlan" class="modal-title">
			Update {{ selectedPlan.name }}
		</div>
		<div v-else class="modal-title">
			{{ selectedPlan.name }}
		</div>
		<div class="modal-content">
			<n-form ref="updatePlanFormRef" :model="selectedPlan" :rules="planRules" :disabled="viewPlan">
				<n-form-item v-if="!viewPlan" label="Recipes" path="recipes" required>
					<n-select max-tag-count="responsive" multiple v-model:value="selectedPlanRecipes" @update:value="updateSelectedPlanRecipes" filterable :options="availableRecipes" />
				</n-form-item>
				<n-form-item v-if="selectedPlanRecipes.length > 0" :label="viewPlan ? 'Recipes' : 'Recipes Quantities'" required>
					<n-space vertical style="width: 100%; max-height: 300px; overflow-y: auto">
						<n-form-item v-for="(recipe, index) in selectedPlan.recipes" :key="index" :label="recipe.name" :show-feedback="false" required>
							<n-input-number :show-buttom="false" :default-value="1" min=1 :precision="0" step=1 v-model:value="recipe.quantity" style="width: 100%" />
						</n-form-item>
					</n-space>
				</n-form-item>
				<n-form-item label="Description" path="description">
					<n-input type="textarea" maxlength="1023" show-count v-model:value="selectedPlan.description" placeholder="Description of the plan" />
				</n-form-item>
			</n-form>
		</div>
		<template #action>
			<n-button v-if="!viewPlan" type="primary" @click="updatePlan">
				Update
			</n-button>
			<n-button v-else type="warning" @click="updateModalVisible = false">
				Close
			</n-button>
		</template>
	</n-modal>
</template>

<style lang="scss" scoped>
@use "../assets/main.scss";

.modal-title {
	color: main.$green;
	font-weight: 600;
	font-size: 28px;
	margin-bottom: 10px;
}

.nutrition-header {
	color: main.$green;
	margin-top: 10px;
}

.divider {
	margin-bottom: 20px;
}
</style>
