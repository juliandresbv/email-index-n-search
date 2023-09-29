<template>
  <div class="mx-7">
    <div class="mb-3">
      <h2 class="text-xl font-semibold leading-tight">Current Email</h2>
    </div>
    <div class="overflow-y-auto max-h-[80vh]" v-if="selectedEmail">
      <div class="mb-5">
        <div class="mt-2">
          <div class="font-bold">From:</div>
          <div class="ml-1 mr-1 break-all">{{ selectedEmail.from }}</div>
        </div>
        <div class="mt-2">
          <div class="font-bold">To:</div>
          <div class="ml-1 mr-1 break-all">{{ selectedEmail.to }}</div>
        </div>
        <div class="mt-2">
          <div class="font-bold">Date:</div>
          <div class="ml-1 mr-1 break-all">{{ selectedEmail.date }}</div>
        </div>
      </div>
      <hr class="h-0.5 mx-1.5 bg-gray-300 border-0 rounded" />
      <div class="mt-5">
        <div class="mt-2">
          <div class="font-bold">Subject:</div>
          <div class="ml-1 mr-1 break-all">{{ selectedEmail.subject }}</div>
        </div>
        <div class="mt-2">
          <div class="font-bold">Body:</div>
          <div class="ml-1 mr-1 break-all">{{ selectedEmail.body }}</div>
        </div>
      </div>
    </div>
    <div class="overflow-y-auto max-h-[80vh]" v-else>
      <div class="mt-2 flex justify-center">
        <p class="text-gray-500">Select an email to show its content</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'

import { useEmailsStore } from '../stores'

const emailsStore = useEmailsStore()

const { selectedEmail } = storeToRefs(emailsStore)
</script>
