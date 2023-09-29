import { defineStore } from 'pinia'

export const useComponentsRefsStore = defineStore('components-refs', {
  state: () => {
    const searchTerm = ''
    const limit = 10
    const page = 1

    return {
      searchTerm,
      limit,
      page
    }
  },
  actions: {
    setSearchTerm(searchTerm: string) {
      this.searchTerm = searchTerm
    },
    setLimit(limit: number) {
      this.limit = limit
    },
    setPage(page: number) {
      this.page = page
    },
    resetSearchTerm() {
      this.searchTerm = ''
    },
    resetLimit() {
      this.limit = 10
    },
    resetPage() {
      this.page = 1
    }
  }
})
