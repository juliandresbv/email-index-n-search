<template>
  <div class="flex items-center justify-between">
    <div class="w-[18vw]">
      <div class="flex items-center justify-start">
        <svg
          class="w-6 h-6 mx-2 text-black dark:text-black"
          aria-hidden="true"
          xmlns="http://www.w3.org/2000/svg"
          fill="currentColor"
          viewBox="0 0 20 16"
        >
          <path
            d="m10.036 8.278 9.258-7.79A1.979 1.979 0 0 0 18 0H2A1.987 1.987 0 0 0 .641.541l9.395 7.737Z"
          />
          <path
            d="M11.241 9.817c-.36.275-.801.425-1.255.427-.428 0-.845-.138-1.187-.395L0 2.6V14a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V2.5l-8.759 7.317Z"
          />
        </svg>
        <span class="font-semibold text-3xl">Emails Search App</span>
      </div>
    </div>
    <div class="w-[72vw]">
      <div
        class="border relative flex items-center h-14 rounded-lg focus-within:shadow-md bg-white overflow-hidden"
      >
        <div class="grid place-items-center h-full w-12 text-gray-300">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-6 w-6"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
            />
          </svg>
        </div>
        <input
          class="peer h-full w-full outline-none text-gray-700 pr-2 font-bold text-xl"
          id="search-emails-input"
          name="search-emails-input"
          type="text"
          placeholder="Search emails..."
          @input="searchEmails"
        />
      </div>
    </div>
    <div class="w-[10vw]">
      <div class="flex flex-col items-center justify-end">
        <span class="text-xs xs:text-sm text-gray-900 font-semibold"> Search type </span>
        <div class="inline-flex mt-1">
          <select
            class="text-sm text-gray-800 py-2 px-4 rounded border"
            name="email-search-type"
            id="email-search-type"
            @change="changeSearchType"
          >
            <option value="match">Match</option>
            <option value="matchphrase">Match phrase</option>
            <option value="term">Term</option>
            <option value="querystring">Query string</option>
            <option value="prefix">Prefix</option>
            <option value="wildcard">Wildcard</option>
            <option value="fuzzy">Fuzzy</option>
          </select>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'

import { debounce } from '../utils'
import { makeRequest } from '../services'
import type { EmailsSearchResponse } from '../types'
import { updateEmailsStoreState } from '../stores/utils'
import { useEmailsStore, useComponentsRefsStore } from '../stores'

const emailsStore = useEmailsStore()
const componentsRefsStore = useComponentsRefsStore()

const { page } = storeToRefs(componentsRefsStore)
const { limit } = storeToRefs(componentsRefsStore)
const { searchTerm } = storeToRefs(componentsRefsStore)
const { searchType } = storeToRefs(componentsRefsStore)

const debounceTimeout = 500

const debouncedMakeRequest = debounce(makeRequest<EmailsSearchResponse>, debounceTimeout)

const searchEmails = async (event: any) => {
  const inputSearchTerm = event.target.value.trim()

  if (inputSearchTerm == null || inputSearchTerm.length <= 0) {
    emailsStore.$reset()
    componentsRefsStore.resetSearchTerm()
    componentsRefsStore.resetPage()

    return
  }

  componentsRefsStore.setSearchTerm(inputSearchTerm)
  componentsRefsStore.resetPage()

  const searchEmailsResponseData = await debouncedMakeRequest('POST', '/emails/search', {
    term: searchTerm.value,
    searchType: searchType.value,
    page: page.value,
    limit: limit.value
  })

  updateEmailsStoreState(searchEmailsResponseData, emailsStore)
}

const changeSearchType = async (event: any) => {
  const inputSearchType = event.target.value

  if (inputSearchType == null || inputSearchType.length <= 0) {
    return
  }

  componentsRefsStore.setSearchType(inputSearchType)

  if (searchTerm.value.length <= 0) {
    return
  }

  const searchEmailsResponseData = await debouncedMakeRequest('POST', '/emails/search', {
    term: searchTerm.value,
    searchType: searchType.value,
    page: page.value,
    limit: limit.value
  })

  updateEmailsStoreState(searchEmailsResponseData, emailsStore)
}
</script>
