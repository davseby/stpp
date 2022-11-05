<script setup lang="ts">
import { ref, computed } from 'vue'
import { useAuth } from "@/stores/useAuth"
import { useProduct } from "@/stores/useProduct"
import { useRecipe } from "@/stores/useRecipe"
import axios from 'axios'

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
	NDivider,
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

const createRecipeFormRef = ref<Object>(null);
const updateRecipeFormRef = ref<Object>(null);

const recipeRules = {
	name: {
		message: "Please input a recipe name",
		required: true,
	},
	image_url: {
		message: "Please input a link to the image",
		required: true,
	},
	products: {
		message: "Please select at least 2 products",
		validator(rule, value, callback) {
			return value.length > 1;
		},
	},
	description: {
		message: "Please describe the recipe",
		required: true,
	},
}

const availableProducts = ref<Array<Object>>([]);

const createdRecipe = ref<Recipe>({});
const createdRecipeProducts = ref<Array<Object>>([]);

const selectedRecipe = ref<Recipe>({});
const selectedRecipeProducts = ref<Array<Object>>([]);

const createModalVisible = ref(false);
const updateModalVisible = ref(false);

const notification = useNotification();

useRecipe().retrieveRecipes();

const deleteRecipe = (id: string) => {
	useRecipe().deleteRecipe(id).then(resp => {
		if (resp.error) {
			notification["error"]({
				content: resp.error,
				duration: 2500,
				keepAliveOnHover: true
			})
			return;
		} 

		notification["success"]({
			content: "Succesfully deleted a recipe",
			duration: 2500,
			keepAliveOnHover: true
		})
	});
}

const showCreateModal = () => {
	createdRecipe.value = {
		products: [],
	};

	createdRecipeProducts.value = [];
	for (let i = 0; i < createdRecipe.value.products.length; i++) {
		createdRecipeProducts.value.push(createdRecipe.value.products[i].id);
	}

	availableProducts.value = [];
	for (let i = 0; i < useProduct().products.length; i++) {
		availableProducts.value.push({
			label: useProduct().products[i].name,
			value: useProduct().products[i].id,
		});
	}

	createModalVisible.value = true;
}

const viewRecipe = ref(true);

const showUpdateModal = (recipe: Recipe, view: boolean) => {
	selectedRecipe.value = JSON.parse(JSON.stringify(recipe));

	selectedRecipeProducts.value = [];
	for (let i = 0; i < selectedRecipe.value.products.length; i++) {
		selectedRecipeProducts.value.push(selectedRecipe.value.products[i].id);
	}

	availableProducts.value = [];
	for (let i = 0; i < useProduct().products.length; i++) {
		availableProducts.value.push({
			label: useProduct().products[i].name,
			value: useProduct().products[i].id,
		});
	}

	viewRecipe.value = view;
	updateModalVisible.value = true;
}

const updateCreatedRecipeProducts = (value: string, option: any) => {
	let products = [];
	for (let i = 0; i < option.length; i++) {
		let found = false;
		for (let j = 0; j < createdRecipe.value.products.length; j++) {
			if (option[i].value == createdRecipe.value.products[j].id) {
				products.push(createdRecipe.value.products[j]);
				found = true;
				break;
			}
		}

		if (!found) {
			for(let j = 0; j < useProduct().products.length; j++) {
				if(option[i].value == useProduct().products[j].id) {
					let product = JSON.parse(JSON.stringify(useProduct().products[j]));
					product.quantity = 1;
					products.push(product);
					break;
				}
			}
		}
	}

	createdRecipe.value.products = products;
}

const updateSelectedRecipeProducts = (value: string, option: any) => {
	let products = [];
	for (let i = 0; i < option.length; i++) {
		let found = false;
		for (let j = 0; j < selectedRecipe.value.products.length; j++) {
			if (option[i].value == selectedRecipe.value.products[j].id) {
				products.push(selectedRecipe.value.products[j]);
				found = true;
				break;
			}
		}

		if (!found) {
			for(let j = 0; j < useProduct().products.length; j++) {
				if(option[i].value == useProduct().products[j].id) {
					let product = JSON.parse(JSON.stringify(useProduct().products[j]));
					product.quantity = 1;
					products.push(product);
					break;
				}
			}
		}
	}

	selectedRecipe.value.products = products;
}

const createRecipe = () => {
	createRecipeFormRef.value.validate((errors) => {
		if (errors) {
			return
		}

		let products = [];
		for (let i = 0; i < createdRecipe.value.products.length; i++) {
			products.push({
				product_id: createdRecipe.value.products[i].id,
				quantity: createdRecipe.value.products[i].quantity,
			})
		}

		let recipe = JSON.parse(JSON.stringify(createdRecipe.value));
		recipe.products = products;
		useRecipe().createRecipe(recipe).then(resp => {
			if (resp.error) {
				notification["error"]({
					content: resp.error,
					duration: 2500,
					keepAliveOnHover: true
				})
				return;
			} 

			notification["success"]({
				content: "Succesfully created a recipe",
				duration: 2500,
				keepAliveOnHover: true
			})

			createModalVisible.value = false;
		});
	}).catch((errors) => {
		console.log(errors);
	})
}

const updateRecipe = () => {
	updateRecipeFormRef.value.validate((errors) => {
		if (errors) {
			return
		}

		let products = [];
		for (let i = 0; i < selectedRecipe.value.products.length; i++) {
			products.push({
				product_id: selectedRecipe.value.products[i].id,
				quantity: selectedRecipe.value.products[i].quantity,
			})
		}

		let recipe = JSON.parse(JSON.stringify(selectedRecipe.value));
		recipe.products = products;

		useRecipe().updateRecipe(recipe).then(resp => {
			if (resp.error) {
				notification["error"]({
					content: resp.error,
					duration: 2500,
					keepAliveOnHover: true
				})
				return;
			} 

			notification["success"]({
				content: "Succesfully updated a recipe",
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

const searchedRecipes = computed(() => {
	return useRecipe().recipes.filter((recipe) => {
		return (
			(!searchSelfCreated.value || recipe.user_id == useAuth().user.id) && recipe.name
				.toLowerCase()
				.indexOf(searchQuery.value.toLowerCase()) != -1
		);
	});
});

</script>

<template>
	<n-space justify="space-between" align="center">
		<div class="page-title">Recipes</div>
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
		<n-grid-item v-for="(recipe, index) in searchedRecipes" :key="index" style="display: flex">
			<n-card :title="recipe.name" hoverable embedded>
				<template #cover>
					<img 
						:src="recipe.image_url == '' ? 'https://media.istockphoto.com/vectors/thumbnail-image-vector-graphic-vector-id1147544807?k=20&m=1147544807&s=612x612&w=0&h=pBhz1dkwsCMq37Udtp9sfxbjaMl27JUapoyYpQm0anc=' : recipe.image_url"
					>
				</template>
				<n-ellipsis line-clamp="3">
					{{ recipe.description }}
				</n-ellipsis>
				<div>
					<n-ellipsis line-clamp="5">
						<div v-for="(product, index) in recipe.products" :key="index"> 
							<div class="nutrition-header">{{ product.name }}</div>
							<div>{{ product.quantity }} {{ product.serving.type }}</div>
						</div>
					</n-ellipsis>
				</div>
				<template #footer>
					<n-space justify="space-between" align="center">
						<n-button :focusable="false" @click="showUpdateModal(recipe, true)" text>
							<template #icon>
								<n-icon size="20px" color="#36ad6a">
									<view-icon />
								</n-icon>
							</template>
							View
						</n-button>
						<n-button v-if="useAuth().authorized && useAuth().user.id == recipe.user_id" :focusable="false" @click="showUpdateModal(recipe, false)" text>
							<template #icon>
								<n-icon size="20px" color="#36ad6a">
									<update-icon />
								</n-icon>
							</template>
							Update
						</n-button>
						<div v-else></div>
						<n-popconfirm v-if="useAuth().authorized && (useAuth().user.admin || useAuth().user.id == recipe.user_id)"  width="310px" placement="top-end" :positive-text="null" negative-text="Confirm" @negative-click="deleteRecipe(recipe.id)">
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
							Are you sure you want to delete the recipe?
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
			Create a new recipe 
		</div>
		<div class="modal-content">
			<n-space justify="center">
				<div class="image-wrapper">
					<n-image
						width="250"
						object-fit="scale-down"
						:src="createdRecipe.image_url"
						fallback-src="https://media.istockphoto.com/vectors/thumbnail-image-vector-graphic-vector-id1147544807?k=20&m=1147544807&s=612x612&w=0&h=pBhz1dkwsCMq37Udtp9sfxbjaMl27JUapoyYpQm0anc="
					/>
				</div>
			</n-space>
			<div class="divider"></div>
			<n-form ref="createRecipeFormRef" :model="createdRecipe" :rules="recipeRules">
				<n-form-item label="Name" path="name" required>
					<n-input type="text" maxlength="255" show-count v-model:value="createdRecipe.name" placeholder="Name of the recipe" />
				</n-form-item>
				<n-form-item label="Image URL" path="image_url" required>
					<n-input type="text" maxlength="1023" show-count v-model:value="createdRecipe.image_url" placeholder="Image representation of the recipe" />
				</n-form-item>
				<n-form-item label="Products" path="products" required>
					<n-select max-tag-count="responsive" multiple v-model:value="createdRecipeProducts" @update:value="updateCreatedRecipeProducts" filterable placeholder="How it is measured in the recipes" :options="availableProducts" />
				</n-form-item>
				<n-form-item v-if="createdRecipeProducts.length > 0" label="Products Quantities" required>
					<n-space vertical style="width: 100%; max-height: 300px; overflow-y: auto">
						<n-form-item v-for="(product, index) in createdRecipe.products" :key="index" :label="product.name" :show-feedback="false" required>
							<n-input-number :show-buttom="false" :default-value="1" min=0.001 :precision="3" v-model:value="product.quantity" style="width: 100%">
							 <template #suffix>
								{{ product.serving.type }}
							</template>
							</n-input-number>
						</n-form-item>
					</n-space>
				</n-form-item>
				<n-form-item label="Description" path="description">
					<n-input type="textarea" maxlength="1023" show-count v-model:value="createdRecipe.description" placeholder="Description of the recipe" />
				</n-form-item>
			</n-form>
		</div>
		<template #action>
			<n-button type="primary" @click="createRecipe">
				Create
			</n-button>
		</template>
	</n-modal>
	<n-modal
		v-model:show="updateModalVisible" 
		preset="dialog" 
		:show-icon="false"
	>
		<div class="modal-title">
			Update {{ selectedRecipe.name }}
		</div>
		<div class="modal-content">
			<n-space justify="center">
				<div class="image-wrapper">
					<n-image
						width="250"
						object-fit="scale-down"
						:src="selectedRecipe.image_url"
						fallback-src="https://media.istockphoto.com/vectors/thumbnail-image-vector-graphic-vector-id1147544807?k=20&m=1147544807&s=612x612&w=0&h=pBhz1dkwsCMq37Udtp9sfxbjaMl27JUapoyYpQm0anc="
					/>
				</div>
			</n-space>
			<div class="divider"></div>
			<n-form ref="updateRecipeFormRef" :model="selectedRecipe" :rules="recipeRules" :disabled="viewRecipe">
				<n-form-item label="Image URL" path="image_url" required>
					<n-input type="text" maxlength="1023" show-count v-model:value="selectedRecipe.image_url" placeholder="Image representation of the recipe" />
				</n-form-item>
				<n-form-item v-if="!viewRecipe" label="Products" path="products" required>
					<n-select max-tag-count="responsive" multiple v-model:value="selectedRecipeProducts" @update:value="updateSelectedRecipeProducts" filterable placeholder="How it is measured in the recipes" :options="availableProducts" />
				</n-form-item>
				<n-form-item v-if="selectedRecipeProducts.length > 0" :label="viewRecipe ? 'Products' : 'Products Quantities'" required>
					<n-space vertical style="width: 100%; max-height: 300px; overflow-y: auto">
						<n-form-item v-for="(product, index) in selectedRecipe.products" :key="index" :label="product.name" :show-feedback="false" required>
							<n-input-number :show-buttom="false" :default-value="1" min=0.001 :precision="3" v-model:value="product.quantity" style="width: 100%">
							 <template #suffix>
								{{ product.serving.type }}
							</template>
							</n-input-number>
						</n-form-item>
					</n-space>
				</n-form-item>
				<n-form-item label="Description" path="description">
					<n-input type="textarea" maxlength="1023" show-count v-model:value="selectedRecipe.description" placeholder="Description of the recipe" />
				</n-form-item>
			</n-form>
		</div>
		<template #action>
			<n-button v-if="!viewRecipe" type="primary" @click="updateRecipe">
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
}

.modal-content {
	.image-wrapper {
		margin: 20px 0;
		border: 1px solid main.$grey;
	}
}

.nutrition-header {
	color: main.$green;
	margin-top: 10px;
}

.divider {
	margin-bottom: 20px;
}
</style>
