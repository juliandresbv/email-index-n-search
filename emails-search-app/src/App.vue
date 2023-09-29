<template>
  <div>
    <div class="mt-4 mb-10 mx-5">
      <SearchBar />
    </div>
    <div class="flex max-h-[90vh]">
      <EmailTable class="w-1/3" />
      <EmailContent class="w-2/3" />
    </div>
  </div>
</template>

<script lang="ts">
import SearchBar from './components/SearchBar.vue'
import EmailTable from './components/EmailTable.vue'
import EmailContent from './components/EmailContent.vue'

export default {
  components: {
    SearchBar,
    EmailTable,
    EmailContent
  }
}
</script>
