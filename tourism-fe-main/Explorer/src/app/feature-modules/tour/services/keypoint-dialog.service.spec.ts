import { TestBed } from '@angular/core/testing';

import { KeypointDialogService } from './keypoint-dialog.service';

describe('KeypointDialogService', () => {
  let service: KeypointDialogService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(KeypointDialogService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
