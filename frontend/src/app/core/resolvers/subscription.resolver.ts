import { inject } from '@angular/core';
import { ResolveFn } from '@angular/router';
import { SubscriptionPlan } from '../models';
import { SubscriptionService } from '../services/subscription.service';

// Ensures subscription data is loaded before any protected route activates.
// Components can then call subscriptionService.hasFeature() synchronously
// without waiting for an async load, eliminating upgrade-banner flash.
export const subscriptionResolver: ResolveFn<SubscriptionPlan | null> = () => {
  return inject(SubscriptionService).load();
};
