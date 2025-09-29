import { TestBed } from '@angular/core/testing';

import { KeypointService } from './keypoint.service';

describe('KeypointService', () => {
  let service: KeypointService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(KeypointService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
