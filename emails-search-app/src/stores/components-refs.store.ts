import { defineStore } from 'pinia'

export const useComponentsRefsStore = defineStore('components-refs', {
  state: () => {
    const searchType = 'match'
    const searchTerm = ''
    const limit = 10
    const page = 1

    return {
      searchType,
      searchTerm,
      limit,
      page
    }
  },
  actions: {
    setSearchType(searchType: string) {
      this.searchType = searchType
    },
    setSearchTerm(searchTerm: string) {
      this.searchTerm = searchTerm
    },
    setLimit(limit: number) {
      this.limit = limit
    },
    setPage(page: number) {
      this.page = page
    },
    resetSearchType() {
      this.searchType = 'match'
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
