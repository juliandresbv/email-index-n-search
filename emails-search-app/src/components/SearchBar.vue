<template>
  <div class="flex items-center">
    <div class="w-1/6">
      <div class="flex items-center">
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
    <div class="w-5/6">
      <div
        class="border relative flex items-center h-12 rounded-lg focus-within:shadow-md bg-white overflow-hidden"
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
          class="peer h-full w-full outline-none text-gray-700 pr-2 font-semibold"
          id="search-emails-input"
          name="search-emails-input"
          type="text"
          placeholder="Search emails..."
          @input="emailsSearchRequest"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'

import { makeRequest } from '../services'
import type { EmailsSearchResponse } from '../types'
import { useEmailsStore, useComponentsRefsStore } from '../stores'

const emailsStore = useEmailsStore()
const componentsRefsStore = useComponentsRefsStore()

const { page } = storeToRefs(componentsRefsStore)
const { limit } = storeToRefs(componentsRefsStore)

const debounceTimeout = 500

let timeoutId: NodeJS.Timeout | null = null

const emailsSearchRequest = async (event: any) => {
  const term = event.target.value.trim()

  if (term.length <= 0) {
    emailsStore.$reset()
    componentsRefsStore.resetSearchTerm()

    return
  }

  if (timeoutId) {
    clearTimeout(timeoutId)
  }

  timeoutId = setTimeout(async () => {
    componentsRefsStore.setSearchTerm(term)
    componentsRefsStore.resetPage()

    const searchEmailsReq = await makeRequest<EmailsSearchResponse>('POST', '/emails/search', {
      term: term,
      page: page.value,
      limit: limit.value
    })

    if (!searchEmailsReq) {
      emailsStore.$reset()

      return
    }

    const { emails, hits } = searchEmailsReq.data

    if (!emails || emails?.length <= 0 || !hits || hits <= 0) {
      emailsStore.$reset()

      return
    }

    emailsStore.setEmails(emails)
    emailsStore.setHits(hits)
  }, debounceTimeout)
}
</script>
